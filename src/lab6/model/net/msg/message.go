package msg

import (
	"encoding/gob"
	"lab6/controller/node"
	"time"
)

var (
	ClientChan chan ClientRequestMessage
)

func init() {
	gob.Register(Node{})
	gob.Register(HARTBEATREQUEST{})
	gob.Register(HARTBEATRESPONSE{})
	gob.Register(Lead{})
	gob.Register(LEADERREQUEST{})
	gob.Register(LEADERRESPONSE{})
	gob.Register(LISTRESPONSE{})
	gob.Register(MACHINECOUNT{})
	gob.Register(MESSAGE{})
	gob.Register(TIMEOUT{})
	gob.Register(ClientRequestMessage{})
	gob.Register(ClientResponseMessage{})
	gob.Register(ClientRequestNodes{})
	gob.Register(ClientResponseNodes{})
	ClientChan = make(chan ClientRequestMessage, 10)
}

type Node struct {
	IP        string
	ROLE      string
	TIME      int64
	SUSPECTED bool
	ALIVE     bool
	LEAD      bool
}

type Lead struct {
	IP        string
	TIME      int64
	ALIVE     bool
	SUSPECTED bool
}

type HARTBEATREQUEST struct {
	IP string
}

type HARTBEATRESPONSE struct {
	IP string
}

type LEADERREQUEST struct {
	TONODE   node.T_Node
	FROMNODE node.T_Node
}

type LEADERRESPONSE struct {
	NODE node.T_Node
}

type LISTRESPONSE struct {
	LIST []node.T_Node
}
type MACHINECOUNT struct {
	I    int
	NODE []node.T_Node
}
type MESSAGE struct {
	MSG string
}

type TIMEOUT struct {
	DELAY time.Duration
}

type ClientRequestMessage struct {
	Content       interface{}
	RemoteAddress string
}

type ClientResponseMessage struct {
	Round     int
	MsgNumber int
	Value     string
}

type ClientRequestNodes struct {
	RemoteAddress string
}

type ClientResponseNodes struct {
	List []node.T_Node
}
