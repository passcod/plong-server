package main

import (
    "fmt"
    "encoding/json"
    "net/http"
    "os"
    "time"
//    "code.google.com/p/go.net/websocket"
    "plong"
)

const Version string = "0.2.6"

type JustOK struct {
  Ok bool
}

func main() {
  plong.Configure(plong.Config{1800})
  
  mux := http.NewServeMux()
  mux.HandleFunc("/", hello)
  mux.HandleFunc("/wuu2", status)
  mux.HandleFunc("/ohai", new_client)
  mux.HandleFunc("/obai", del_client)
  mux.HandleFunc("/iam", new_identity)
  mux.HandleFunc("/whois", find_identity)
  
  port := os.Getenv("PORT")
  if port == "" {
    port = "1501"
  }
  
  serv := &http.Server{
    Addr: ":" + port,
    Handler: mux,
    ReadTimeout: 30 * time.Second,
    WriteTimeout: 30 * time.Second,
  }
  
  fmt.Printf("Plong server v.%s started.\n", Version)
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

func set_headers(res http.ResponseWriter, code int) {
  res.Header().Set("Access-Control-Allow-Origin", "*")
  res.Header().Set("Content-type", "application/json")
  res.WriteHeader(code)
}


func hello(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    set_headers(res, 200)
    
    type Hello struct {
      Version string
      Http bool
      WebSocket bool
    }
    
    enc := json.NewEncoder(res)
    hi := Hello{Version, true, false}
    
    if err := enc.Encode(&hi); err != nil {
      fmt.Println(err)
    }
}


func status(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    set_headers(res, 200)
    
    type Status struct {
      Peers int
      Identities int
    }
    
    enc := json.NewEncoder(res)
    enc.Encode(Status{plong.PeerCount(), plong.IdentityCount()})
}

func new_client(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    set_headers(res, 200)
    
    client := plong.NewPeer()
    
    enc := json.NewEncoder(res)
    enc.Encode(client)
}

func del_client(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    enc := json.NewEncoder(res)
    dec := json.NewDecoder(req.Body)

    type jString struct {
      Id string
    }
    
    var j jString
    if err := dec.Decode(&j); err != nil {
      set_headers(res, 400)
      enc.Encode(JustOK{false})
      fmt.Println(err)
      return
    }

    peer := plong.FindPrivatePeer(j.Id)
    peer.Destroy()
    
    set_headers(res, 200)
    enc.Encode(JustOK{true})
}

func new_identity(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    
    enc := json.NewEncoder(res)
    dec := json.NewDecoder(req.Body)

    type IdReq struct {
      Private string
      Passphrase string
    }
    
    var ir IdReq
    if err := dec.Decode(&ir); err != nil || ir.Private == "" || ir.Passphrase == "" {
      set_headers(res, 400)
      enc.Encode(JustOK{false})
      fmt.Println(err)
      return
    }
    
    peer := plong.FindPrivatePeer(ir.Private)
    peer.NewIdentity(ir.Passphrase)
    
    set_headers(res, 200)
    enc.Encode(JustOK{true})
}

func find_identity(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    
    enc := json.NewEncoder(res)
    dec := json.NewDecoder(req.Body)

    type jPass struct {
      Passphrase string
    }
    
    var j jPass
    if err := dec.Decode(&j); err != nil {
      set_headers(res, 400)
      enc.Encode(JustOK{false})
      fmt.Println(err)
      return
    }
    
    i, ok := plong.FindIdentity(j.Passphrase)
    if !ok {
      set_headers(res, 404)
      enc.Encode(JustOK{false})
      return
    }
    
    type IdResp struct {
      Public string
      Created time.Time
    }
    
    set_headers(res, 200)
    enc.Encode(IdResp{i.Subject.PublicId, i.CreatedAt})
}