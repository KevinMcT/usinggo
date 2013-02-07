package messages

import (
	"encoding/gob"
)

// Types use by gob as interface{} have to be registered
func init() {
	gob.Register(Ping{})
}

/* Structs for the messages: */

type Ping struct {
	IP   string
	Time float64
	Role string
}
