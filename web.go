package main

import (
    "fmt"
    "encoding/json"
    "net/http"
    "os"
    "bytes"
    "math/rand"
    "strconv"
)

var VERSION [2]string = [2]string{"0.0.10", "Binary Shite"}
const MAX_CLIENTS int = 100000 // <— arbitrary

type Clients map[int]int
var clients Clients = make(Clients)

type JustOK struct {
  Ok bool
}


func main() {
    http.HandleFunc("/", hello)
    http.HandleFunc("/wuu2", status)
    http.HandleFunc("/ohai", new_client)
    http.HandleFunc("/obai", del_client)
    http.HandleFunc("/iam", new_identity)
    http.HandleFunc("/whois", find_identity)
    http.HandleFunc("/talk", talk_to)
    http.HandleFunc("/hear", hear_from)
    
    fmt.Printf("Plong server v.%s “%s” started.\n", VERSION[0], VERSION[1])
    fmt.Printf("Listening on port %s...\n", os.Getenv("PORT"))
    err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
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

func find_private(public int) int {
  for priv, pub := range clients {
    if (pub == public) {
      return priv
    }
  }
  return -1
}

func hello(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    
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

func status(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    
    type Status struct {
      ConnectedClients int
      MaxClients int
    }
    
    enc := json.NewEncoder(res)
    enc.Encode(Status{len(clients), MAX_CLIENTS})
}

func new_client(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    
    enc := json.NewEncoder(res)
    
    if len(clients) == MAX_CLIENTS {
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
    
    
    type ClientIDs struct {
      Private int
      Public int
    }
    enc.Encode(ClientIDs{new_priv, new_pub})
}

func del_client(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    
    enc := json.NewEncoder(res)
    
    buf := new(bytes.Buffer)
    buf.ReadFrom(req.Body)
    body := buf.String()
    
    id, err := strconv.Atoi(body)
    if err != nil {
      enc.Encode(JustOK{false})
      return
    }
    
    _, ok := clients[id]
    if !ok {
      enc.Encode(JustOK{false})
      return
    }
    
    delete(clients, id)
    
    enc.Encode(JustOK{true})
}

func new_identity(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    
    fmt.Fprintln(res, "new identity")
}

func find_identity(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    
    fmt.Fprintln(res, "find identity")
}

func talk_to(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    
    fmt.Fprintln(res, "talk")
}

func hear_from(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    
    fmt.Fprintln(res, "hear")
}
