package main

import (
	"plong"
)

type Hub map[plong.Peer][]*Connection
var wsHub Hub = make(Hub)


func (h Hub) Direct(c *Connection) int {
	cs, ok := h[c.link]
	if !ok {
		return -1
	}
	
	for i, con := range cs {
		if con.link == c.peer {
			return i
		}
	}
	
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
