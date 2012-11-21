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

// A million clients should be plenty.
// This is not completely arbirtrary:
// I estimate ~50 potential bytes per
// client, thus 50MiB max memory usage
// (although we'll probably hit an HTTP
// bottleneck before reaching that).
const MAX_CLIENTS int = 1000000
var clientIDs []int

func main() {
    http.HandleFunc("/", hello)
    http.HandleFunc("/wuu2", status)
    http.HandleFunc("/ohai", new_client)
    http.HandleFunc("/obai", del_client)
    http.HandleFunc("/iam", new_identity)
    http.HandleFunc("/whois", find_identity)
    http.HandleFunc("/talk", talk_to)
    http.HandleFunc("/hear", hear_from)
    
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

func indexOf(slice []int, value int) int {
  for p, v := range slice {
    if (v == value) {
      return p
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
    hi := Hello{"0.0.7", true, false}
    
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
    enc.Encode(Status{len(clientIDs), MAX_CLIENTS})
}

func new_client(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    
    if len(clientIDs) == MAX_CLIENTS {
      fmt.Fprintf(res, "%d", -1)
      return
    }
       
    new_id := rand.Int()
    for indexOf(clientIDs, new_id) != -1 {
      new_id = rand.Int()
    }
    clientIDs = append(clientIDs, new_id)
    
    fmt.Fprintf(res, "%d", new_id)
}

func del_client(res http.ResponseWriter, req *http.Request) {
    log_request(req)
    
    type RespOK struct {
      Ok bool
    }
    
    enc := json.NewEncoder(res)
    
    buf := new(bytes.Buffer)
    buf.ReadFrom(req.Body)
    body := buf.String()
    
    id, err := strconv.Atoi(body)
    if err != nil {
      enc.Encode(RespOK{false})
      return
    }
    
    i := indexOf(clientIDs, id)
    if i == -1 {
      enc.Encode(RespOK{false})
      return
    }
    
    // Remove from slice
    clientIDs = append(clientIDs[:i], clientIDs[i+1:]...)
    
    enc.Encode(RespOK{true})
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
