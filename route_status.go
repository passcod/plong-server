package main

import (
	"net/http"
	"plong"
)

func routeStatus(res http.ResponseWriter, req *http.Request) {
	log_request(req)

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

	respond(res, 200, hello{Version, Mode, stats{plong.PeerCount(), plong.IdentityCount()}, PlongConfig.IdentityTimeout})
}
