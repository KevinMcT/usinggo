package RoundVar

import (
	"lab5/model/Network/message"
)

type roundSingle struct {
	Round         int
	MessageNumber int
	List          []message.Node
}

var instantiated *roundSingle = nil

func GetRound() *roundSingle {
	if instantiated == nil {
		instantiated = new(roundSingle)
		instantiated.Round = 0
		instantiated.List = make([]message.Node, 0)
	}
	return instantiated
}
