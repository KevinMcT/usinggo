package messages

import (
	"encoding/gob"
)

// Types use by gob as interface{} have to be registered
func init() {
	gob.Register(StrMsg{})
	gob.Register(ErrMsg{})
	gob.Register(NT{})
}

/* Structs for the messages: */

type StrMsg struct {
	Sender  string
	Content string
}

type ErrMsg struct {
	Sender string
	Error  string
}

type NT struct {
	Sender string
}
