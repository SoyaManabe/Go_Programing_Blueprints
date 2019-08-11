package main
import(
  "log"
  "net/http"
  "github.com/gorilla/websocket"
)

type room struct {
  forward chan []byte
  join chan *client
  leave chan *client
  clients map[*client]bool
}

// newRoom returns new chatroom
func newRoom() *room {
  return &room{
    forward: make(chan []byte),
    join: make(chan *client),
    leave: make(chan *client),
    clients: make(map[*client]bool),
  }
}

func (r *room) run() {
  for{
    select {
    case client := <-r.join:
      // join
      r.clients[client] = true
    case client := <-r.leave:
      // leave
      delete(r.clients, client)
      close(client.send)
    case msg := <-r.forward:
      // send msg for all clients
      for client := range r.clients {
        select {
        case client.send <- msg:
          // send msg
        default:
          // failue
          delete(r.clients, client)
          close(client.send)
        }
      }
    }
  }
}

const (
  socketBufferSize = 1024
  messageBufferSize = 256
)
var upgrader = &websocket.Upgrader{ReadBufferSize:
    socketBufferSize, WriteBufferSize: socketBufferSize}
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  socket, err := upgrader.Upgrade(w, req, nil)
  if err != nil {
    log.Fatal("serveHTTP:", err)
    return
  }
  client := &client{
    socket: socket,
    send: make(chan []byte, messageBufferSize),
    room: r,
  }
  r.join <- client
  defer func() { r.leave <- client }()
  go client.write()
  client.read()
}