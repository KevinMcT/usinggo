package Discov

import (
	"fmt"
	"lab3/messages"
	"net"
	"strings"
	"time"
)

var (
	port   = 2100
	ticker = time.NewTicker(time.Second * 5)
)

func init() {

}

func ListenForBroadcast(inChan chan messages.Node, intime string, leadChan chan *net.UDPAddr) {
	fmt.Println("LISTEN FOR BROADCAST")
	mcaddr, _ := net.ResolveUDPAddr("udp4", "239.255.43.99:1888")
	conn, _ := net.ListenMulticastUDP("udp4", nil, mcaddr)
	go SendBroadcast(mcaddr, conn, intime, leadChan)
	for {
		data := make([]byte, 512)
		n, addr, _ := conn.ReadFromUDP(data)
		readFrom := string(data[0:n])
		res := strings.Split(readFrom, ":")
		if addr != nil {
			node := messages.Node{IP: addr, Time: res[0], Lead: res[1]}
			inChan <- node
		}
		<-ticker.C
	}
}

func SendBroadcast(mcaddr *net.UDPAddr, conn *net.UDPConn, intime string, leadChan chan *net.UDPAddr) {
	timer := time.NewTicker(5 * time.Second)
	for {
		fmt.Println("Hanging here....")
		lead := <-leadChan
		conn.WriteTo([]byte(intime+":"+lead.IP.String()), mcaddr)
		<-timer.C
	}
}
