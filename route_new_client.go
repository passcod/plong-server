package main

import (
	"net/http"
	"plong"
)

func routeNewClient(res http.ResponseWriter, req *http.Request) {
	log_request(req)

	respond(res, 200, plong.NewPeer())
}
