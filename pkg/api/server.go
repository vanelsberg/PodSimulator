package api

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
    "github.com/avereha/pod/pkg/pod"
)

type Server struct {
	pod   *pod.Pod
}

func New(pod *pod.Pod) *Server {

  ret := &Server{
    pod:   pod,
  }

  return ret
}

func (s *Server) Start() {
    fmt.Println("Pod simulator web api listening")
    s.setupRoutes()
    http.ListenAndServe(":8080", nil)
}

func (s *Server) setupRoutes() {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Simple Server")
  })
  // map our `/ws` endpoint to the `serveWs` function
  http.HandleFunc("/ws", serveWs)
}

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
  ReadBufferSize:  1024,
  WriteBufferSize: 1024,

  // We'll need to check the origin of our connection
  // this will allow us to make requests from our React
  // development server to here.
  // For now, we'll do no checking and just allow any connection
  CheckOrigin: func(r *http.Request) bool { return true },
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn) {
  for {
  // read in a message
    messageType, p, err := conn.ReadMessage()
    if err != nil {
      log.Println(err)
      return
    }
// print out that message for clarity
    fmt.Println(string(p))

    if err := conn.WriteMessage(messageType, p); err != nil {
      log.Println(err)
      return
    }
  }
}

// define our WebSocket endpoint
func serveWs(w http.ResponseWriter, r *http.Request) {
  fmt.Println(r.Host)

  // upgrade this connection to a WebSocket
  // connection
  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println(err)
  }
  // listen indefinitely for new messages coming
  // through on our WebSocket connection
  reader(ws)
}
