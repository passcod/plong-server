package main

import (
	"encoding/json"
	"net/http"
	"plong"
	"time"
)

func routeFindIdentity(res http.ResponseWriter, req *http.Request) {
	log_request(req)

	dec := json.NewDecoder(req.Body)

	type pass struct {
		Passphrase string
	}

	var j pass
	if err := dec.Decode(&j); err != nil {
		respond(res, 400, err)
		return
	}

	i, ok := plong.FindIdentity(j.Passphrase)
	if !ok {
		respond(res, 404, nil)
		return
	}

	type identityResponse struct {
		Public  string
		Created time.Time
	}

	respond(res, 200, identityResponse{i.Subject.PublicId, i.CreatedAt})
}
