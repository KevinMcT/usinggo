package main

import (
	"fmt"
	//"lab3/Communication"
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
	//ok, ip := Discov.Listener()
	startTime := time.Now()
	//ticker := time.NewTicker(7 * time.Second)
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

/*	var ok bool
	ok = true
	for {
		fmt.Println("REG IP")
		node := <-ips
		if ok {
			ipArray = append(ipArray, node.IP)
			ok = false
		}
		fmt.Println(ipArray)
		for _, v := range ipArray {
			if node.IP != v {
				ipArray = append(ipArray, node.IP)
				fmt.Println(ipArray)
			}
		}
		work <- 1
	}*/
