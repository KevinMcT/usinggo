package LeaderElect

import (
	"fmt"
	"lab3Test/model/Network/message"
)

var ()

func Elect(leadElect chan []message.Node) {
	list := <-leadElect
	fmt.Println("In leader Elect")
	fmt.Println("NodeList: ", list)
	fmt.Println("Oldest process:::")
	old := findOldest(list)
	fmt.Println(old)
}

func findOldest(list []message.Node) message.Node {
	var node message.Node
	for _, v := range list {
		if node.IP == "" {
			node = v
		} else {
			if v.TIME < node.TIME && v.SUSPECTED == false {
				node = v
			}
		}
	}
	return node
}
