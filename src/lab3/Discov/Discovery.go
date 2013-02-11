package Discov

import (
	"fmt"
	"lab3/messages"
	"net"
	"time"
)

var (
	port   = 2100
	ticker = time.NewTicker(time.Second * 5)
)

func init() {

}

func ListenForBroadcast(inChan chan messages.Node, intime string) {
	fmt.Println("LISTEN FOR BROADCAST")
	mcaddr, _ := net.ResolveUDPAddr("udp4", "239.255.43.99:1888")
	conn, _ := net.ListenMulticastUDP("udp4", nil, mcaddr)
	go SendBroadcast(mcaddr, conn, intime)
	for {
		data := make([]byte, 512)
		n, addr, _ := conn.ReadFromUDP(data)
		hisTime := string(data[0:n])
		if addr != nil {
			node := messages.Node{IP: addr, Time: hisTime}
			inChan <- node
		}
		<-ticker.C
	}
}

func SendBroadcast(mcaddr *net.UDPAddr, conn *net.UDPConn, intime string) {
	timer := time.NewTicker(5 * time.Second)
	for {
		conn.WriteTo([]byte(intime), mcaddr)
		<-timer.C
	}
}
