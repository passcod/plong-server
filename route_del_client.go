package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"plong"
)

func routeDelClient(res http.ResponseWriter, req *http.Request) {
	log_request(req)
	dec := json.NewDecoder(req.Body)

	type anid struct {
		Id string
	}

	var j anid
	if err := dec.Decode(&j); err != nil {
		respond(res, 400, err)
		fmt.Println(err)
		return
	}

	peer := plong.FindPrivatePeer(j.Id)
	peer.Destroy()

	respond(res, 200, nil)
}
