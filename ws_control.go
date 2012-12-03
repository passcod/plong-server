package main

import (
	"fmt"
	"strings"
	"encoding/json"
)

func (c *Connection) Control(message string) {
	fmt.Printf("[Debug] Control sequence received: ‘%s’.\n", message)
	split := strings.Split(message, " ")
	command := split[0]
	args := split[1:]
	
	switch command {
	case "ping":
		c.send <- "\x1bpong"
	case "wuu2":
		wsStatus(c, args)
/*	case "iam":
		wsNewIdentity(c, args)
	case "whois":
		wsFindIdentity(c, args) */
	case "dalink":
		wsLinkStatus(c, args)
	case "chlink":
		wsLinkChange(c, args)
	}
}

func wsJson(c *Connection, v interface{}) {
	type fail struct {
		Error string
	}

	if v != nil {
		b, err := json.Marshal(v)
		if err != nil {
			fmt.Println(err)
			b, err = json.Marshal(fail{err.Error()})
		}
		c.send <- fmt.Sprintf("\x1b%s", b)
	}
}