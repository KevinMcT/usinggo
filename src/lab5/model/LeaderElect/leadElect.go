package LeaderElect

import (
	"lab5/model/Network/message"
)

var ()

func Elect(leadElect chan []message.Node, elected chan message.Node, block chan int) {
	list := <-leadElect
	old := findOldest(list)
	elected <- old
	block <- 1
}

func findOldest(list []message.Node) message.Node {
	var node message.Node
	for i, v := range list {
		if node.IP == "" {
			node = v
		} else {
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
