package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"plong"
	"strings"
)

type Connection struct {
	ws   *websocket.Conn
	send chan string
}

var wsHubConn map[*Connection]plong.Peer = make(map[*Connection]plong.Peer)
var wsHubPeer map[plong.Peer]*Connection = make(map[plong.Peer]*Connection)

func wsHandler(ws *websocket.Conn) {
	c := &Connection{send: make(chan string, BufferSize), ws: ws}
	url := strings.Split(fmt.Sprintf("%s", ws.LocalAddr()), "/")
	id := url[len(url)-1]
	peer := plong.FindPublicPeer(id)

	if peer.PrivateId == "" {
		fmt.Printf("[WebSocket] Error: No such peer “%s”.\n", id)
		return
	}

	wsHubConn[c] = peer
	wsHubPeer[peer] = c

	fmt.Printf("[WebSocket] New connection: “%s”.\n", id)
}
