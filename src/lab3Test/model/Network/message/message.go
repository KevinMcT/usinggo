package message

import (
	"encoding/gob"
)

var ()

func init() {
	gob.Register(Node{})
}

type Node struct {
	IP        string
	TIME      int64
	LEAD      bool
	ALIVE     bool
	SUSPECTED bool
}
