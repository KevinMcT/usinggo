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
}

func Listen(nodeChan chan message.Node, startTime int64, exit chan bool, nLead message.Node) {
	name1, _ := os.Hostname()
	addr1, _ := net.LookupHost(name1)
	UDPAddr1, _ := net.ResolveUDPAddr("udp4", addr1[0]+":1888")
	if nLead.IP == "" {
		mcaddr1, _ := net.ResolveUDPAddr("udp4", "239.255.43.99:1888")
		conn1, _ := net.ListenMulticastUDP("udp4", nil, mcaddr1)
		conn1.SetDeadline(time.Now().Add(5 * time.Second))
		data1 := make([]byte, 1024)
		_, _, err1 := conn1.ReadFromUDP(data1)
		if err1 != nil {
			fmt.Println("No leader, setting self")
			leader = message.Node{IP: UDPAddr1.IP.String(), ALIVE: true, LEAD: true, SUSPECTED: false}
			leaderChan = true
			fst = true
		} else {
			leaderChan = false
		}
		conn1.Close()
	} else {
		if UDPAddr1.IP.String() == nLead.IP {
			fmt.Println("Leader is selected from leader elect")
			nLead.LEAD = true
			leader = nLead
			fmt.Println("New LEADER UDP", leader)
			leaderChan = true
			fst = true
		} else {
			leaderChan = false
		}
	}

	//---------------------------------------------------------------------------------

	mcaddr, _ := net.ResolveUDPAddr("udp4", "239.255.43.99:1888")
	conn, _ := net.ListenMulticastUDP("udp4", nil, mcaddr)
	if fst {
		nodeChan <- leader
		fst = false
	}
	if leaderChan {
		go BroadcastLeader(mcaddr, conn, leader, startTime, exit)
	} else {
		go Broadcast(mcaddr, conn, startTime, exit)
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
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(100 * time.Millisecond)
			timeout <- true
		}()
		select {
		case <-exit:
			fmt.Println("Break udp listen")
			break
		case <-timeout:
		}
	}
}

func Broadcast(mcaddr *net.UDPAddr, conn *net.UDPConn, startTime int64, exit chan bool) {
	timer := time.NewTicker(1000 * time.Millisecond)
	for {
		conn.WriteTo([]byte(strconv.FormatInt(startTime, 10)), mcaddr)
		<-timer.C
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(100 * time.Millisecond)
			timeout <- true
		}()
		select {
		case <-exit:
			fmt.Println("Break udp broadcast")
			break
		case <-timeout:
		}
	}
}

func BroadcastLeader(mcaddr *net.UDPAddr, conn *net.UDPConn, lead message.Node, startTime int64, exit chan bool) {
	timer := time.NewTicker(1000 * time.Millisecond)
	for {
		str := fmt.Sprintf("lead:%s:%s", lead.IP, strconv.FormatInt(startTime, 10))
		conn.WriteTo([]byte(str), mcaddr)
		<-timer.C
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(100 * time.Millisecond)
			timeout <- true
		}()
		select {
		case <-exit:
			fmt.Println("Break udp broadcastLeader")
			break
		case <-timeout:
		}
	}
}
