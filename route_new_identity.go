package main

import (
	"encoding/json"
	"net/http"
	"plong"
)

func routeNewIdentity(res http.ResponseWriter, req *http.Request) {
	log_request(req)

	dec := json.NewDecoder(req.Body)

	type identRequest struct {
		Private    string
		Passphrase string
	}

	var ir identRequest
	if err := dec.Decode(&ir); err != nil {
		respond(res, 400, err)
		return
	}

	if ir.Private == "" || ir.Passphrase == "" {
		respond(res, 400, "There’s a blank field or two here.")
		return
	}

	peer := plong.FindPrivatePeer(ir.Private)
	if peer.PublicId == "" {
		respond(res, 404, "No such peer.")
		return
	}
	
	peer.NewIdentity(ir.Passphrase)

	respond(res, 200, nil)
}
