package main

import (
	"lab3Test/model/Network/message"
	"lab3Test/model/Network/udp"
)

var (
	nodeChan = make(chan message.Node, 10)
	nodeList = make([]message.Node, 0)
	work     = make(chan int, 0)
)

func main() {
	//startTime := time.Now().Unix()
	go udp.Listen(nodeChan)
	go RegIP()
	for {
		<-work
	}
}

func RegIP() {
	for {
		node := <-nodeChan
		nodeList = AppendIfMissing(nodeList, node)
		work <- 1
	}
}

func AppendIfMissing(slice []message.Node, i message.Node) []message.Node {
	for _, ele := range slice {
		if ele.IP == i.IP {
			return slice
		}
	}
	return append(slice, i)
}
