package message

import (
	"encoding/gob"
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
	LearnChan = make(chan Wrapper, 10)
	AcceptChan = make(chan Wrapper, 10)
	PrepareChan = make(chan Wrapper, 10)
	PromiseChan = make(chan Wrapper, 10)
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
