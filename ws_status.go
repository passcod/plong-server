package main

import (
	"plong"
)

func wsStatus(c *Connection, args []string) {
	type stats struct {
		Peers      int
		Identities int
	}

	type hello struct {
		Version         string
		Mode            string
		Status          stats
		IdentityTimeout int64
	}

	wsJson(c, hello{Version, Mode, stats{plong.PeerCount(), plong.IdentityCount()}, PlongConfig.IdentityTimeout}, false)
}