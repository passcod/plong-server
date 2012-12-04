Plong 
=====

_Pseudo-peer-to-peer API._

COMPILE
-------

I assume you have a working [Go](http://golang.org) install, set up using the
recommended `~/go/{bin,pkg,src}/` structure.

1. Clone this repo into ~/go/src/plong: `git clone git://github.com/passcod/plong-server.git ~/go/src/plong`.
2. Cd into it: `cd ~/go/src/plong`.
3. Switch to plong-lib branch: `git checkout plong-lib`.
4. Build & install the lib: `go build && go install`.
5. Clone this repo (yes, again) into ~/go/src/plong-server: `git clone git://github.com/passcod/plong-server.git ~/go/src/plong-server`.
6. Cd into it: `cd ~/go/src/plong-server`.
7. Get dependencies & build the server: `go get`.


INSTALL
-------

There's binaries for a few systems in the [Dowloads](https://github.com/passcod/plong-server/downloads) section.
These are not really tested nor updated very often.


SYNOPSIS
--------

    $ plong-server
	Plong server v.1.0.0 started (mode: hwx).
    Listening on port 1501...


CONFIG
------

In `config.json`:

 - identity_timeout: How long before an identity dies.
 - buffer_size: The websocket/passthru buffer size.
 - mode: Which features are enabled. This is notably useful to disable WebSockets (`w`) or passthru (`x`).


USAGE
-----

This describes a basic 2-player session. The clients are named ‘Bob’ and ‘Alice’
and there is only one Plong server. You can try it for yourself using the demo
servers at http://canna.plong.me:1501 and http://lotus.plong.me:1501

1. Bob connects to the server.
   
   ```javascript
   $.get('http://server:1501/ohai', function (peer) {
     window.peer = peer
     // This looks like:
     // {
     //   PrivateId: "a long string",
     //   PublicId: "a long string",
     //   Created: "a date/time string"
     // }
   }, 'json')
   ```

2. Bob creates an identity.
   
   ```javascript
   $.ajax({
     type: 'POST',
     url: 'http://server:1501/iam',
     data: JSON.stringify({
       PrivateId: peer.PrivateId,
       Passphrase: "something short and sweet"
     }),
     dataType: 'json',
     processData: false
   });
   ```

3. Alice connects to the server.
   
   ```javascript
   // From another tab/window/browser
   $.get('http://server:1501/ohai', function (peer) {
     window.peer = peer
   }, 'json')
   ```

4. Out of band, Bob gives Alice the passphrase for his identity.
   This must be done within 30 minutes as the identity expires
   after that time.

5. Alice retrieves Bob's PublicId using the passphrase.
   
   ```javascript
   $.ajax({
     type: 'POST',
     url: 'http://server:1501/whois',
     data: JSON.stringify({
       Passphrase: "something short and sweet"
     }),
     dataType: 'json',
     processData: false,
     success: function (bob) {
       window.bob = bob
       // This looks like:
       // {
       //   PublicId: "bob's public id",
       //   Created: "when bob connected"
       // }
     }
   });
   ```

6. Now Alice opens a WebSocket connection:
   
   ```javascript
   // Note we use our own PrivateId to connect.
   ws = new WebSocket("ws://server:1501/ws/"+peer.PrivateId);
   
   // Wait for the connection to open…
   ...
   
   // Now open a link to Bob:
   ws.send("\x1bchlink "+bob.PublicId);
   ```

7. Meanwhile, Bob has also opened a WebSocket connection, and awaits a message:
   
   ```javascript
   ws = new WebSocket("ws://server:1501/ws/"+peer.PrivateId);
   
   // Wait for the connection to open…
   ...
   
   // Listen for messages:
   ws.onmessage = function (m) { ... }
   ```

8. Alice sends a message to Bob:
   
   ```javascript
   ws.send("Hi");
   ```

9. Bob doesn't have a direct link to Alice yet (he cannot
   open it because he doesn't have Alice's Public Id). So
   he receives this message:
   
   ```plain
   :Au7dk...(Alice's Public Id)...8sHsbf=:Hi
   ```

10. Now Bob has Alice's Public Id, he opens a direct link:
   
    ```javascript
    ws.send("\x1bchlink " + alice_public_id)
    ```

11. Now both Alice and Bob have open links to each other,
    when they communicate they will not have the `:Id:`
    header.


PROTOCOL
--------

See the source (`route_*.go` for HTTP, `ws_*.go` for WebSocket) for the moment,
I'll write thorough documentation later on.


LICENSE
-------

No license! This is in the [Public Domain](http://passcod.net/license.html)! :)