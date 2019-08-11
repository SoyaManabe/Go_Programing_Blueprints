package main
import (
  "github.com/gorilla/websocket"
  "time"
)
//client shows one user using chat app
type client struct {
  // socket is Websocket for this client
  socket *websocket.Conn
  // send is a channel send message
  send chan *message
  //room is a room client joining
  room *room
  // userData obtains information about User
  userData map[string]interface{}
 }

func (c *client) read() {
  for {
    var msg *message
    if err := c.socket.ReadJSON(&msg); err == nil {
      msg.When = time.Now()
      msg.Name = c.userData["name"].(string)
      c.room.forward <- msg
    }else{
      break
    }
  }
  c.socket.Close()
}

func (c *client) write() {
  for msg := range c.send {
    if err := c.socket.WriteJSON(msg);
        err != nil {
      break
    }
  }
  c.socket.Close()
}
