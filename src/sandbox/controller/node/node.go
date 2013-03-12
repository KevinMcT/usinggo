package node

import (
	"sandbox/controller/machine"
)

var ()

type T_Node struct {
	IP        string
	ROLE      string
	TIME      int64
	SUSPECTED bool
	ALIVE     bool
	LEAD      bool
}

func CreateNode(comp machine.T_Machine) T_Node {
	var node T_Node
	node.IP = comp.IP
	node.TIME = comp.TIME
	node.ALIVE = true
	node.SUSPECTED = false
	return node
}
