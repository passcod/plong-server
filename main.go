package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"plong"
	"time"
)

func main() {
	readConfig()
	plong.Configure(PlongConfig)

	mux := http.NewServeMux()
	
	// Always handle root, even if the HTTP API
	// is disabled.
	mux.HandleFunc("/", routeStatus)
	
	if whatMode("h") {
		mux.HandleFunc("/wuu2", routeStatus)
		mux.HandleFunc("/ohai", routeNewClient)
		mux.HandleFunc("/obai", routeDelClient)
		mux.HandleFunc("/iam", routeNewIdentity)
		mux.HandleFunc("/whois", routeFindIdentity)
	}
	
	if whatMode("w") {
		mux.Handle("/ws/", websocket.Handler(wsHandler))
	}
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "1501"
	}

	serv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("Plong server v.%s started (mode: %s).\n", Version, Mode)
	fmt.Printf("Listening on port %s...\n", port)
	err := serv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func log_request(req *http.Request) {
	fmt.Printf("%s --> %s %s %s (%s)\n",
		req.RemoteAddr,
		req.Proto,
		req.Method,
		req.RequestURI,
		req.Header["User-Agent"][0])
}

// Sets the proper headers and encodes the value provided to JSON.
func respond(res http.ResponseWriter, code int, v interface{}) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-type", "application/json")
	res.WriteHeader(code)

	type fail struct {
		Error string
	}

	if code != 200 {
		v = fail{fmt.Sprint(v)}
		fmt.Println(v)
	}

	if v != nil {
		enc := json.NewEncoder(res)
		if err := enc.Encode(&v); err != nil {
			fmt.Println(err)
			enc.Encode(fail{err.Error()})
		}
	}
}
