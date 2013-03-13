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
		if node.IP == "" && v.SUSPECTED == false {
			node = v
		} else {
			if v.LEAD == true && v.SUSPECTED == true {
				list[i].LEAD = false
			}
			if v.TIME < node.TIME && v.SUSPECTED == false {
				fmt.Println(list[i])
				list[i].LEAD = true
				node = list[i]
			}
		}
	}
	return node
}
