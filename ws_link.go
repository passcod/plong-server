package main

import (
	"fmt"
	"plong"
)

func (c *Connection) PassThru(message string) {
	fmt.Printf("[Debug] Message received: ‘%s’.\n", message)
	if c.link.PublicId != "" {
		fmt.Printf("[WebSocket] Sending message to ‘%s’.\n", c.link.PublicId)
		
		type fail struct {
			Error string
		}
		
		rem, ok := wsHub[c.link]
		if !ok {
			fmt.Printf("%#v\n", rem)
			wsJson(c, fail{"Peer not connected."})
			return
		}
		
		remote := rem[0]
		direct := wsHub.Direct(c.peer, c.link)
		if direct != -1 {
			remote = rem[direct]
		} else {
			message = fmt.Sprintf(":%s:%s", c.link.PublicId, message)
		}
		remote.send <- message
	}
}

func wsLinkStatus(c *Connection, args []string) {
	type linkstat struct {
		Remote interface{}
		Direct bool
	}
	
	if c.link.PublicId != "" {
		wsJson(c, linkstat{c.link.PublicId, wsHub.Direct(c.peer, c.link) != -1})
	} else {
		wsJson(c, linkstat{Remote: false})
	}
}

func wsLinkChange(c *Connection, args []string) {
	type fail struct {
		Error string
	}
	
	peer := plong.FindPublicPeer(args[0])
	if peer.PublicId == "" {
		wsJson(c, fail{"No such peer."})
		return
	}
	
	fmt.Printf("[WebSocket] Changing link to ‘%s’.\n", peer.PublicId)
	c.link = peer
}