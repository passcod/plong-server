package main

import (
    "fmt"
    "encoding/json"
    "net/http"
    "os"
    "bytes"
    "math/rand"
    "strconv"
    "time"
)

var VERSION [2]string = [2]string{"0.1.0", "Internet Truffle"}
var MAX_CLIENTS int = 100000 // <— arbitrary
var IDENT_TIMEOUT int64 = 1800

type Ident struct {
  Private int
  Created time.Time
}

type JustOK struct {
  Ok bool
}

var clients map[int]int = make(map[int]int)
var idents map[string]Ident = make(map[string]Ident)


func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", hello)
  mux.HandleFunc("/wuu2", status)
  mux.HandleFunc("/ohai", new_client)
  mux.HandleFunc("/obai", del_client)
  mux.HandleFunc("/iam", new_identity)
  mux.HandleFunc("/whois", find_identity)
  
  serv := &http.Server{
    Addr: ":" + os.Getenv("PORT"),
    Handler: mux,
    ReadTimeout: 30 * time.Second,
    WriteTimeout: 30 * time.Second,
  }
  
  fmt.Printf("Plong server v.%s “%s” started.\n", VERSION[0], VERSION[1])
  fmt.Printf("Listening on port %s...\n", os.Getenv("PORT"))
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
    hi := Hello{VERSION[0], true, false}
    
    if err := enc.Encode(&hi); err != nil {
      fmt.Println(err)
    }
}


func check_idents() int {
  for k, i := range idents {
    if _, ok := clients[i.Private]; !ok || i.Created.Before(time.Unix(time.Now().Unix() - IDENT_TIMEOUT, 0)) {
      delete(idents, k)
    }
  }
  
  return len(idents)
}

func status(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    set_headers(res, 200)
    
    type Status struct {
      ConnectedClients int
      MaxClients int
      ActiveIdents int
    }
    
    enc := json.NewEncoder(res)
    enc.Encode(Status{len(clients), MAX_CLIENTS, check_idents()})
}

func new_client(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    
    enc := json.NewEncoder(res)
    
    if len(clients) == MAX_CLIENTS {
      set_headers(res, 503)
      enc.Encode(JustOK{false})
      return
    }
       
    new_pub, new_priv := rand.Int(), rand.Int()
    
    // Make sure new_priv is unique
    for {
      _, ok := clients[new_priv]
      if !ok {
        break
      }
      new_priv = rand.Int()
    }
    
    // Make sure new_pub is unique
    for _, pub := range clients {
      for pub == new_pub {
        new_pub = rand.Int()
      }
    }
    
    // Actually create it
    clients[new_priv] = new_pub
    
    
    type ClientID struct {
      Private int
      Public int
    }
    
    set_headers(res, 200)
    enc.Encode(ClientID{new_priv, new_pub})
}

func del_client(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    
    enc := json.NewEncoder(res)
    
    buf := new(bytes.Buffer)
    buf.ReadFrom(req.Body)
    body := buf.String()
    
    id, err := strconv.Atoi(body)
    if err != nil {
      set_headers(res, 400)
      enc.Encode(JustOK{false})
      return
    }
    
    _, ok := clients[id]
    if !ok {
      set_headers(res, 404)
      enc.Encode(JustOK{false})
      return
    }
    
    delete(clients, id)
    
    set_headers(res, 200)
    enc.Encode(JustOK{true})
}

func new_identity(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    
    enc := json.NewEncoder(res)
    dec := json.NewDecoder(req.Body)
    
    type IdReq struct {
      Private int
      Passphrase string
    }
    
    var ir IdReq
    if err := dec.Decode(&ir); err != nil {
      set_headers(res, 400)
      enc.Encode(JustOK{false})
      return
    }
    
    id := Ident{ir.Private, time.Now()}
    idents[ir.Passphrase] = id
    
    set_headers(res, 200)
    enc.Encode(JustOK{true})
}

func find_identity(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    
    enc := json.NewEncoder(res)
    
    buf := new(bytes.Buffer)
    buf.ReadFrom(req.Body)
    phrase := buf.String()
    
    check_idents()
    
    i, ok := idents[phrase]
    if !ok {
      set_headers(res, 404)
      enc.Encode(JustOK{false})
      return
    }
    
    type IdResp struct {
      Public int
      Created time.Time
    }
    
    set_headers(res, 200)
    enc.Encode(IdResp{clients[i.Private], i.Created})
}