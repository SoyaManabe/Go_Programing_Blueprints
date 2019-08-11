package main
import (
  "time"
)
// message describes one message
type message struct {
  Name string
  Message string
  When time.Time
}
