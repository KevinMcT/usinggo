package message

import (
	"encoding/gob"
	"fmt"
)

var (
	LearnChan   chan Wrapper
	AcceptChan  chan Wrapper
	PrepareChan chan Wrapper
	PromiseChan chan Wrapper
)

func init() {
	gob.Register(Prepare{})
	gob.Register(Promise{})
	gob.Register(Accept{})
	gob.Register(Learn{})
	fmt.Println("Creating chans")
	LearnChan = make(chan Wrapper, 10)
	AcceptChan = make(chan Wrapper, 10)
	PrepareChan = make(chan Wrapper, 10)
	PromiseChan = make(chan Wrapper, 10)
	fmt.Println("done crating chans")
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
	LASTACCEPTEDVALUE string
}

type Accept struct {
	ROUND int
	VALUE string
}

type Learn struct {
	ROUND int
	VALUE string
}
