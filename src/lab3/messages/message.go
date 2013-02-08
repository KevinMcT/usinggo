package messages

import (
	"encoding/gob"
	"net"
	"time"
)

// Types use by gob as interface{} have to be registered
func init() {
	gob.Register(Ping{})
	gob.Register(Node{})
}

/* Structs for the messages: */

type Ping struct {
	Time time.Time
}

type Node struct {
	IP   *net.UDPAddr
	Time time.Time
}
