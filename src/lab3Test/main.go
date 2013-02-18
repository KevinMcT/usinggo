package main

import (
	"fmt"
	"lab3Test/model/FailureDetect"
	"lab3Test/model/LeaderElect"
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
	leadElect = make(chan []message.Node, 10)
	elected   = make(chan message.Node, 1)
	work      = make(chan int, 0)
	wait      = make(chan int, 0)
	leader    message.Node
	nLead     message.Node
	tick      = time.NewTimer(5 * time.Second)
	selfnode  message.Node
	exitUdp   = make(chan bool, 0)
	exitReg   = make(chan bool, 0)
)

func main() {
	startTime := time.Now().UnixNano()
	name, _ := os.Hostname()
	addr, _ := net.LookupHost(name)
	UDPAddr, _ := net.ResolveUDPAddr("udp4", addr[0]+":1888")
	go udp.Listen(nodeChan, startTime, exitUdp, nLead)
	go RegIP(exitReg)
	<-tick.C
	time.Sleep(2 * time.Second)
	if leader.IP != "" && leader.IP == UDPAddr.IP.String() {
		selfnode = message.Node{IP: UDPAddr.IP.String(), TIME: startTime, ALIVE: true, LEAD: true}
	}
	if leader.IP != UDPAddr.IP.String() {
		selfnode = message.Node{IP: UDPAddr.IP.String(), TIME: startTime, ALIVE: true, LEAD: false}
	}
	for {
		if !contains(nodeList, true) {
			fmt.Println("Has no leader....")
			exitReg <- true
			nodeList = make([]message.Node, 0)
			go RegIP(exitReg)
		}
		go FailureDetect.Fd(newNodes, selfnode, leadElect)
		go LeaderElect.Elect(leadElect, elected, work)
		<-work
		exitReg <- true
		exitUdp <- true
		exitUdp <- true
		exitUdp <- true
		newLd := <-elected
		if UDPAddr.IP.String() != newLd.IP {
			time.Sleep(5 * time.Second)
		}
		nodeList = make([]message.Node, 0)
		nodeChan = make(chan message.Node, 10)
		newNodes = make(chan message.Node, 10)
		nLead = newLd
		if nLead.IP == UDPAddr.IP.String() {
			fmt.Println("Leader in main")
			selfnode = nLead
			selfnode.TIME = startTime
			selfnode.ALIVE = true
			selfnode.LEAD = true
			selfnode.IP = UDPAddr.IP.String()
			nLead.LEAD = true
		}
		if nLead.IP != UDPAddr.IP.String() {
			fmt.Println("Slave in main")
			selfnode = message.Node{IP: UDPAddr.IP.String(), TIME: startTime, ALIVE: true, LEAD: false}
		}
		go udp.Listen(nodeChan, startTime, exitUdp, nLead)
		go RegIP(exitReg)
	}
}

func RegIP(exit chan bool) {
	for {
		node := <-nodeChan
		nodeList = AppendIfMissing(nodeList, node)
		for _, v := range nodeList {
			if v.LEAD == true {
				leader = v
				break
			}
		}
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(100 * time.Millisecond)
			timeout <- true
		}()
		select {
		case <-exit:
			fmt.Println("Break RegIP")
			break
		case <-timeout:
		}
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

func contains(s []message.Node, e bool) bool {
	for _, a := range s {
		if a.LEAD == e {
			return true
		}
	}
	return false
}
