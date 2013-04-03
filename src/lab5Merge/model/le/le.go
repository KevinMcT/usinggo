package le

import (
	"lab5Merge/controller/node"
)

var ()

func Elect(leadElect []node.T_Node) node.T_Node {
	list := leadElect
	old := findOldest(list)
	return old
}

func findOldest(list []node.T_Node) node.T_Node {
	var node node.T_Node
	for i, v := range list {
		if node.IP == "" && v.SUSPECTED == false {
			node = v
			node.LEAD = true
		} else if node.IP != "" && v.SUSPECTED == false {
			if v.LEAD == true && v.SUSPECTED == true {
				list[i].LEAD = false
			}
			if v.TIME < node.TIME && v.SUSPECTED == false {
				list[i].LEAD = true
				node = list[i]
			}
		}
	}
	return node
}
