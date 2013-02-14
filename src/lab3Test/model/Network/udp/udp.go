package udp

import (
	"fmt"
	"lab3Test/model/Network/message"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	port       = 2100
	ticker     = time.NewTicker(time.Millisecond * 100)
	leaderChan bool
	fst        bool
	leader     message.Node
)

func init() {
	name, _ := os.Hostname()
	addr, _ := net.LookupHost(name)
	UDPAddr, _ := net.ResolveUDPAddr("udp4", addr[0]+":1888")
	mcaddr, _ := net.ResolveUDPAddr("udp4", "239.255.43.99:1888")
	conn, _ := net.ListenMulticastUDP("udp4", nil, mcaddr)
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	data := make([]byte, 1024)
	_, _, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println("No leader, setting self")
		leader = message.Node{IP: UDPAddr.IP.String(), ALIVE: true, LEAD: true, SUSPECTED: false}
		leaderChan = true
		fst = true
	} else {
		leaderChan = false
	}
}

func Listen(nodeChan chan message.Node, startTime int64) {
	mcaddr, _ := net.ResolveUDPAddr("udp4", "239.255.43.99:1888")
	conn, _ := net.ListenMulticastUDP("udp4", nil, mcaddr)
	if fst {
		nodeChan <- leader
		fst = false
	}
	if leaderChan {
		go BroadcastLeader(mcaddr, conn, leader, startTime)
	} else {
		go Broadcast(mcaddr, conn, startTime)
	}
	for {
		data := make([]byte, 1024)
		n, addr, _ := conn.ReadFromUDP(data)
		recived := string(data[0:n])
		recivedSplit := strings.Split(recived, ":")
		if addr != nil {
			if strings.Contains(recivedSplit[0], "lead") {
				t, _ := strconv.ParseInt(recivedSplit[2], 10, 64)
				node := message.Node{IP: addr.IP.String(), TIME: t, ALIVE: true, LEAD: true, SUSPECTED: false}
				nodeChan <- node
			} else {
				t, _ := strconv.ParseInt(recived, 10, 64)
				node := message.Node{IP: addr.IP.String(), TIME: t, ALIVE: true, LEAD: false, SUSPECTED: false}
				nodeChan <- node
			}
		}
		<-ticker.C
	}
}

func Broadcast(mcaddr *net.UDPAddr, conn *net.UDPConn, startTime int64) {
	timer := time.NewTicker(100 * time.Millisecond)
	for {
		conn.WriteTo([]byte(strconv.FormatInt(startTime, 10)), mcaddr)
		<-timer.C
	}
}

func BroadcastLeader(mcaddr *net.UDPAddr, conn *net.UDPConn, lead message.Node, startTime int64) {
	timer := time.NewTicker(100 * time.Millisecond)
	for {
		conn.WriteTo([]byte("lead:"+lead.IP+":"+strconv.FormatInt(startTime, 10)), mcaddr)
		<-timer.C
	}
}
