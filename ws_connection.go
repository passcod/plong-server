package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"plong"
	"strings"
)

type Hub map[plong.Peer][]*Connection
var wsHub Hub = make(Hub)


func (h Hub) Direct(local plong.Peer, remote plong.Peer) int {
	return -1
}

func (h Hub) Add(c *Connection) {
	_, ok := h[c.peer]
	if !ok {
		cs := []*Connection{c}
		h[c.peer] = cs
	} else {
		h[c.peer] = append(h[c.peer], c)
	}
}

func (h Hub) Remove(c *Connection) {
	cs, ok := h[c.peer]
	if ok {
		ids := []int{}
		w := 0
		for i, con := range cs {
			if c == con {
				ids = append(ids, i)
			}
		}
		
loop:
		for i, x := range cs {
			for _, id := range ids {
				if id == i {
					continue loop
				}
			}
			cs[w] = x
			w++
		}
		
		h[c.peer] = cs[:w]
	}
}

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
	wsHub.Add(c)
	fmt.Printf("[WebSocket] New connection: “%s”.\n", id)

	defer func() {
		wsHub.Remove(c)
	}()

	go c.Writer()
	c.Reader()
}
