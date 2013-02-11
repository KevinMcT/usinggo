package main

import (
	"fmt"
	"lab3/Discov"
	"lab3/LE"
	"lab3/messages"
	"net"
	"os"
	"strconv"
	"time"
)

var (
	end     = make(chan int)
	ips     = make(chan messages.Node, 100)
	work    = make(chan int, 10)
	ipArray = make([]messages.Node, 0)
	Leader  *net.UDPAddr
)

func main() {
	startTime := time.Now().Unix()
	test := strconv.FormatInt(startTime, 10)
	name, _ := os.Hostname()
	addr, _ := net.LookupHost(name)
	UDPAddr, _ := net.ResolveUDPAddr("udp4", addr[0]+":1888")
	node := messages.Node{IP: UDPAddr, Time: test}
	ipArray = append(ipArray, node)
	go Discov.ListenForBroadcast(ips, test)
	go RegIP()
	for {
		<-work
		if Leader == nil {
			ip := LeaderElect.Elect(test, ipArray)
			if ip == nil {
				Leader = UDPAddr
			} else {
				Leader = ip
			}
		}
		fmt.Println(ipArray)
	}
}

func RegIP() {
	for {
		node := <-ips
		ipArray = AppendIfMissing(ipArray, node)
		work <- 1
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
