package main

import (
	"plong"
	"strings"
	"time"
	"encoding/json"
)

func wsNewIdentity(c *Connection, args []string) {
	arg := []byte(strings.Join(args, " "))
	type identRequest struct {
		PrivateId  string
		Passphrase string
	}

	var ir identRequest
	if err := json.Unmarshal(arg, &ir); err != nil {
		wsJson(c, err, true)
		return
	}

	if ir.PrivateId == "" || ir.Passphrase == "" {
		wsJson(c, "There’s a blank field or two here.", true)
		return
	}

	peer := plong.FindPrivatePeer(ir.PrivateId)
	if peer.PublicId == "" {
		wsJson(c, "No such peer.", true)
		return
	}
	
	peer.NewIdentity(ir.Passphrase)
}

func wsFindIdentity(c *Connection, args []string) {
	arg := []byte(strings.Join(args, " "))
	type pass struct {
		Passphrase string
	}

	var j pass
	if err := json.Unmarshal(arg, &j); err != nil {
		wsJson(c, err, true)
		return
	}

	i, ok := plong.FindIdentity(j.Passphrase)
	if !ok {
		wsJson(c, "Identity not found.", true)
		return
	}

	type identityResponse struct {
		PublicId string
		Created  time.Time
	}

	wsJson(c, identityResponse{i.Subject.PublicId, i.CreatedAt}, false)
}
