package le

import (
	"lab6/controller/node"
)

/*
	Will search all the nodes in the nodeList and select the oldest node
	which is not suspected.
	The age of the node is determined by the timestamp sent in the UDP Broadcast	
*/

/*
	Public function
	leadElect 	[]node.T_Node	 The systems node list
	returns the oldest node
*/
func Elect(leadElect []node.T_Node) node.T_Node {
	list := leadElect
	old := findOldest(list)
	return old
}

/*
	Private function
	list 	[]node.T_Node	 The list to be searched
	returns the oldest node
*/
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
