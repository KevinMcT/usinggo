package leaderelection

import (
	"fmt"
	"sandbox/controller/node"
)

var ()

func Elect(leadElect []node.T_Node, elected chan node.T_Node, block chan int) node.T_Node {
	list := leadElect
	old := findOldest(list)
	return old
	//	elected <- old
	//	block <- 1
}

func findOldest(list []node.T_Node) node.T_Node {
	var node node.T_Node
	for i, v := range list {
		if node.IP == "" {
			node = v
			node.LEAD = true
		} else {
			if v.TIME < node.TIME && v.SUSPECTED == false {
				fmt.Println("Should find leader")
				list[i].LEAD = true
				node = list[i]
			}
		}
	}
	return node
}
