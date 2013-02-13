package udp

import (
	"fmt"
	"lab3Test/model/Network/message"
	"net"
	"os"
	"time"
)

var (
	port       = 2100
	ticker     = time.NewTicker(time.Second * 1)
	leaderChan bool
	leader     message.Node
)

func init() {
	name, _ := os.Hostname()
	addr, _ := net.LookupHost(name)
	UDPAddr, _ := net.ResolveUDPAddr("udp4", addr[0]+":1888")
	mcaddr, _ := net.ResolveUDPAddr("udp4", "239.255.43.99:1888")
	conn, _ := net.ListenMulticastUDP("udp4", nil, mcaddr)
	conn.SetDeadline(time.Now().Add(1 * time.Second))
	data := make([]byte, 512)
	_, _, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println("No leader, setting self")
		leader = message.Node{IP: UDPAddr.IP.String(), ALIVE: true, LEAD: true, SUSPECTED: false}
		leaderChan = true
	} else {
		leaderChan = false
	}
}

func Listen(nodeChan chan message.Node) {
	mcaddr, _ := net.ResolveUDPAddr("udp4", "239.255.43.99:1888")
	conn, _ := net.ListenMulticastUDP("udp4", nil, mcaddr)
	if leaderChan {
		go BroadcastLeader(mcaddr, conn, leader)
	} else {
		go Broadcast(mcaddr, conn)
	}
	for {
		data := make([]byte, 512)
		n, addr, _ := conn.ReadFromUDP(data)
		if addr != nil {
			fmt.Println(string(data[0:n]))
			node := message.Node{IP: addr.IP.String(), ALIVE: true, LEAD: false}
			nodeChan <- node
		}
		<-ticker.C
	}
}

func Broadcast(mcaddr *net.UDPAddr, conn *net.UDPConn) {
	timer := time.NewTicker(1 * time.Second)
	for {
		conn.WriteTo([]byte(""), mcaddr)
		<-timer.C
	}
}

func BroadcastLeader(mcaddr *net.UDPAddr, conn *net.UDPConn, lead message.Node) {
	timer := time.NewTicker(1 * time.Second)
	for {
		conn.WriteTo([]byte("lead:"+lead.IP), mcaddr)
		<-timer.C
	}
}
