package plong

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

// An Identity is a time-limited token
// used to retrieve info about a Peer.
type Identity struct {
	Subject   Peer
	Phrase    string
	CreatedAt time.Time
}

// A Peer is a pair of unique IDs:
// the PrivateId is used by the client
// to communicate with the server, and
// the PublicId is used by all others
// to communicate with the client (thru
// the server).
type Peer struct {
	PrivateId string
	PublicId  string
	CreatedAt time.Time
}

// Used to set configurable options.
type Config struct {
	IdentityTimeout int64
}

// The PeerList is mapped by Id. There's
// two lists set up, one mapped by Private
// and the other by Public. This is mostly
// for speed and easy of access, but it could
// be done using a single slice instead.
type PeerList map[string]Peer

// The IdentityStore is mapped by token.
// There's no duplication here, as the
// only way to get the associated Peer
// should be to provide the token. But
// Identities expire, so we still need
// to cycle through occasionally to
// delete a bunch of them. *mmmmm*
type IdentityStore map[string]Identity

/// Globals ///
var publicList, privateList PeerList
var identities IdentityStore
var config Config

// Set the config. This is required before
// anything else can be run.
func Configure(newConfig Config) {
	config = newConfig
	publicList = make(PeerList)
	privateList = make(PeerList)
	identities = make(IdentityStore)
}

// Clean the store of every invalid Identity.
// An Indentity is invalid if it is older than
// what Config.IdentityTimeout specifies, or
// if its Subject isn't in the Lists anymore.
func (store IdentityStore) Clean() {
	for token, ident := range store {
		_, ok := privateList[ident.Subject.PrivateId]
		thatTime := time.Unix(time.Now().Unix()-config.IdentityTimeout, 0)

		if !ok || ident.CreatedAt.Before(thatTime) {
			delete(store, token)
		}
	}
}

// Create a new client. We create two random IDs, of 256b
// and 512b for the Public and Private respectively,
// and make sure they're unique. Then we add them to
// the Lists. And finally we return the Peer.
func NewPeer() Peer {
	newPub, newPriv := randomString(256), randomString(512)

	for {
		_, ok := privateList[newPriv]
		if !ok {
			break
		}
		newPriv = randomString(512)
	}

	for {
		_, ok := publicList[newPub]
		if !ok {
			break
		}
		newPub = randomString(256)
	}

	newPeer := Peer{newPriv, newPub, time.Now()}

	privateList[newPriv] = newPeer
	publicList[newPub] = newPeer

	return newPeer
}

// Find a peer by its PrivateId
func FindPrivatePeer(id string) Peer {
	p, _ := privateList[id]
	return p
}

// Find a peer by its PublicId
func FindPublicPeer(id string) Peer {
	p, _ := publicList[id]
	return p
}

// Destroy the peer and remove it from the
// lists. (Actually, just the fact of removing
// it from the lists should be enough to
// ensure it gets collected as garbage by Go.)
//
// We also destroy all identities the peer had.
func (peer Peer) Destroy() {
	delete(privateList, peer.PrivateId)
	delete(publicList, peer.PublicId)

	for _, ident := range identities {
		if ident.Subject == peer {
			ident.Destroy()
		}
	}
}

// Create a n-bytes Base64 random string.
// The returned string will be bigger than
// n because it generates n bytes and then
// encodes it.
func randomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return err.Error()
	}

	enc := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01234567890-_")
	return enc.EncodeToString(b)
}

// Create an identity for a peer
func (peer Peer) NewIdentity(phrase string) Identity {
	newIdentity := Identity{peer, phrase, time.Now()}
	identities[phrase] = newIdentity
	return newIdentity
}

// Destroy an identity
func (ident Identity) Destroy() {
	delete(identities, ident.Phrase)
}

// Find an identity
func FindIdentity(phrase string) (Identity, bool) {
	i, r := identities[phrase]
	return i, r
}

// Return the number of peers
func PeerCount() int {
	return len(privateList)
}

// Return the number of identities.
func IdentityCount() int {
	return len(identities)
}
