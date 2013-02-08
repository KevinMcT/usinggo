package main

import (
	"fmt"
	"lab3/Discov"
	"lab3/messages"
	"net"
	"os"
	"time"
)

var (
	end     = make(chan int)
	ips     = make(chan messages.Node, 100)
	work    = make(chan int, 10)
	ipArray = make([]messages.Node, 0)
)

func main() {
	startTime := time.Now()
	name, _ := os.Hostname()
	addr, _ := net.LookupHost(name)
	UDPAddr, _ := net.ResolveUDPAddr("udp4", addr[0]+":1888")
	node := messages.Node{IP: UDPAddr, Time: startTime}
	ipArray = append(ipArray, node)
	go Discov.ListenForBroadcast(ips, startTime)
	go RegIP()
	for {
		<-work
	}
}

func RegIP() {
	for {
		node := <-ips
		ipArray = AppendIfMissing(ipArray, node)
		fmt.Println(ipArray)
	}
}

func AppendIfMissing(slice []messages.Node, i messages.Node) []messages.Node {
	for _, ele := range slice {
		if ele.IP.String() == i.IP.String() {
			return slice
		}
	}
	return append(slice, i)
}
