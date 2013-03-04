package message

import (
	"encoding/gob"
)

var (
	ClientChan chan ClientRequestMessage
)

func init() {
	gob.Register(Node{})
	gob.Register(HARTBEATREQUEST{})
	gob.Register(HARTBEATRESPONSE{})
	gob.Register(ClientRequestMessage{})
	ClientChan = make(chan ClientRequestMessage, 10)
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
