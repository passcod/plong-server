package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"plong"
	"strings"
)

var wsHub map[plong.Peer][]*Connection = make(map[plong.Peer][]*Connection)

type Connection struct {
	ws   *websocket.Conn
	send chan string
	peer plong.Peer
	link plong.Peer
}

func (c *Connection) Reader() {
	for {
		var message string
		err := websocket.Message.Receive(c.ws, &message)
		if err != nil {
			fmt.Printf("[WebSocket] Error: %s.\n", err)
			break
		}

		if message[0] == '\x1b' {
			c.Control(message[1:])
		} else {
			c.PassThru(message)
		}
	}
	c.Close()
}

func (c *Connection) Writer() {
	for message := range c.send {
		err := websocket.Message.Send(c.ws, message)
		if err != nil {
			fmt.Printf("[WebSocket] Error: %s.\n", err)
			break
		}
	}
	c.Close()
}

func (c *Connection) Close() {
	fmt.Printf("[Websocket] Closing connection: “%s”.\n", c.peer.PublicId)

	c.ws.Close()
}

func wsHandler(ws *websocket.Conn) {
	url := strings.Split(fmt.Sprintf("%s", ws.LocalAddr()), "/")
	id := url[len(url)-1]
	peer := plong.FindPublicPeer(id)

	if peer.PrivateId == "" {
		fmt.Printf("[WebSocket] Error: No such peer “%s”.\n", id)
		return
	}

	c := &Connection{send: make(chan string, BufferSize), ws: ws, peer: peer}

	fmt.Printf("[WebSocket] New connection: “%s”.\n", id)

	defer func() {

	}()

	go c.Writer()
	c.Reader()
}
