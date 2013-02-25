package message

import (
	"encoding/gob"
)

var ()

func init() {
	gob.Register(Node{})
	gob.Register(HARTBEATREQUEST{})
	gob.Register(HARTBEATRESPONSE{})
	gob.Register(ClientRequestMessage{})
}

type Node struct {
	IP        string
	TIME      int64
	LEAD      bool
	ALIVE     bool
	SUSPECTED bool
}

type HARTBEATREQUEST struct {
	IP string
}

type HARTBEATRESPONSE struct {
	IP string
}

type ClientRequestMessage struct {
	Content string
}
