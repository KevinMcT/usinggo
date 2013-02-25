package main

import (
	"encoding/gob"
	"fmt"
	"lab4/Utils"
	"lab4/model/FailureDetect"
	"lab4/model/LeaderElect"
	"lab4/model/Network/message"
	"lab4/model/Network/udp"
	"net"
	"os"
	"time"
)

var (
	nodeChan   = make(chan message.Node, 10)
	newNodes   = make(chan message.Node, 10)
	nodeList   = make([]message.Node, 0)
	leadElect  = make(chan []message.Node, 10)
	elected    = make(chan message.Node, 1)
	work       = make(chan int, 0)
	wait       = make(chan int, 0)
	leader     message.Node //My leader node
	nLead      message.Node
	tick       = time.NewTimer(5 * time.Second)
	clientTick = time.NewTimer(10 * time.Second)
	selfnode   message.Node
	exitUdp    = make(chan bool, 0)
	exitReg    = make(chan bool, 0)
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
	go ClientConnection()
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
			time.Sleep(200 * time.Millisecond)
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

func ClientConnection() {
	fmt.Println("Im waiting bitches!!")
	service := "0.0.0.0:1337"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	Utils.CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	for {
		Utils.CheckError(err)
		conn, err := listener.Accept()

		// I`m leader, must handle request
		if leader.IP == selfnode.IP {
			fmt.Println("Im leader, doing shit!")

			decoder := gob.NewDecoder(conn)
			fmt.Println("Got a message give me a sec...")
			var msg interface{}
			err = decoder.Decode(&msg)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(msg)
			if msg != nil {
				var clientMsg message.ClientRequestMessage
				clientMsg = msg.(message.ClientRequestMessage)
				fmt.Println("Here`s the message:")
				fmt.Println(clientMsg.Content)
				fmt.Println("Im out!")
			}
		} else { // not leader, send it to the leader node
			fmt.Println("Im not leader, send it on!")
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
