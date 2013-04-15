package msg

import (
	"encoding/gob"
	"lab6/model/SlotList"
)

var (
	LearnChan       chan Wrapper
	AcceptChan      chan Wrapper
	PrepareChan     chan Wrapper
	PromiseChan     chan Wrapper
	RestartProposer chan string
	SendPrepareChan chan bool
	GetNodeOnTrack  chan string
)

func init() {
	gob.Register(Prepare{})
	gob.Register(Promise{})
	gob.Register(Accept{})
	gob.Register(Learn{})
	gob.Register(UpdateNode{})
	LearnChan = make(chan Wrapper, 10)
	AcceptChan = make(chan Wrapper, 10)
	PrepareChan = make(chan Wrapper, 10)
	PromiseChan = make(chan Wrapper, 10)
	SendPrepareChan = make(chan bool, 10)
	RestartProposer = make(chan string, 10)
	GetNodeOnTrack = make(chan string, 10)
}

type Wrapper struct {
	Ip      string
	Message interface{}
}

type Prepare struct {
	ROUND int
}

type Promise struct {
	ROUND             int
	LASTACCEPTEDROUND int
	LASTACCEPTEDVALUE interface{}
}

type Accept struct {
	ROUND     int
	MSGNUMBER int
	VALUE     interface{}
}

type Learn struct {
	ROUND     int
	MSGNUMBER int
	VALUE     interface{}
}

type UpdateNode struct {
	PrepareMessage Prepare
	SlotList       *SlotList.Slots
	BankAccounts   map[string]int
}
