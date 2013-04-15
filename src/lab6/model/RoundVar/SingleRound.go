package RoundVar

import (
	"lab5/controller/node"
)

type roundSingle struct {
	Round         int
	MessageNumber int
	RespondClient string
	List          []node.T_Node
	CurrentLeader node.T_Node
}

var instantiated *roundSingle = nil

func GetRound() *roundSingle {
	if instantiated == nil {
		instantiated = new(roundSingle)
		instantiated.Round = 0
		instantiated.List = make([]node.T_Node, 0)
	}
	return instantiated
}
