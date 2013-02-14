package main

import (
	"lab3Test/model/FailureDetect"
	"lab3Test/model/Network/message"
	"lab3Test/model/Network/udp"
	"net"
	"os"
	"time"
)

var (
	nodeChan  = make(chan message.Node, 10)
	newNodes  = make(chan message.Node, 10)
	nodeList  = make([]message.Node, 0)
	leadElect = make(chan int)
	stopNumL  = make(chan int, 10)
	stopNumS  = make(chan int, 10)
	work      = make(chan int, 0)
	leader    message.Node
	tick      = time.NewTimer(5 * time.Second)
	selfnode  message.Node
)

func main() {
	startTime := time.Now().UnixNano()
	name, _ := os.Hostname()
	addr, _ := net.LookupHost(name)
	UDPAddr, _ := net.ResolveUDPAddr("udp4", addr[0]+":1888")
	go udp.Listen(nodeChan, startTime)
	go RegIP()
	<-tick.C
	if leader.IP != "" && leader.IP == UDPAddr.IP.String() {
		selfnode = message.Node{IP: UDPAddr.IP.String(), TIME: startTime, ALIVE: true, LEAD: true}
	}
	if leader.IP != UDPAddr.IP.String() {
		selfnode = message.Node{IP: UDPAddr.IP.String(), TIME: startTime, ALIVE: true, LEAD: false}
	}
	go FailureDetect.Fd(newNodes, selfnode, leadElect)
	for {
		<-work
	}
}

func RegIP() {
	for {
		node := <-nodeChan
		nodeList = AppendIfMissing(nodeList, node)
		for _, v := range nodeList {
			if v.LEAD == true {
				leader = v
				break
			}
		}
		work <- 1
	}
}

func AppendIfMissing(slice []message.Node, i message.Node) []message.Node {
	for _, ele := range slice {
		if ele.IP == i.IP {
			return slice
		}
	}
	newNodes <- i
	return append(slice, i)
}
