package main

import (
	"encoding/gob"
	"fmt"
	"lab5/Utils"
	"lab5/model/FailureDetect"
	"lab5/model/LeaderElect"
	"lab5/model/Network/message"
	"lab5/model/Network/udp"
	"lab5/model/Paxos"
	"lab5/model/RoundVar"
	"net"
	"os"
	"time"
)

var (
	nodeChan      = make(chan message.Node, 10)
	newNodes      = make(chan message.Node, 10)
	newNodesPaxos = make(chan message.Node, 10)

	acceptorChan = make(chan string, 10)
	proposerChan = make(chan message.Learn, 10)
	nodeList     = make([]message.Node, 0)
	leadElect    = make(chan []message.Node, 10)
	elected      = make(chan message.Node, 1)
	work         = make(chan int, 0)
	wait         = make(chan int, 0)
	leader       message.Node //My leader node
	nLead        message.Node
	tick         = time.NewTimer(5 * time.Second)
	clientTick   = time.NewTimer(10 * time.Second)
	selfnode     message.Node
	exitUdp      = make(chan bool, 0)
	exitReg      = make(chan bool, 0)
)

func main() {
	startTime := time.Now().UnixNano()
	name, _ := os.Hostname()
	addr, _ := net.LookupHost(name)
	UDPAddr, _ := net.ResolveUDPAddr("udp4", addr[0]+":1890")
	go udp.Listen(nodeChan, startTime, exitUdp, nLead)
	go RegIP(exitReg)
	go ClientConnection()
	go Paxos.Acceptor()
	go Paxos.Learner()
	go Paxos.PaxosHandler()

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

		go Paxos.Proposer(leader, selfnode, newNodesPaxos, acceptorChan)
		go FailureDetect.Fd(newNodes, selfnode, leadElect)
		go LeaderElect.Elect(leadElect, elected, work)

		//When something comes inn on this channel we want to start the system 
		//over again.
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

/*
Method for listning after client requests. If this instance is the
leader of the paxos system he sends the message on to the proposer. If
not the message is redirected to the node that is the leader.
*/
func ClientConnection() {
	fmt.Println("Waiting for inncoming clients")
	service := "0.0.0.0:1337"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	Utils.CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	for {
		Utils.CheckError(err)
		conn, _ := listener.Accept()
		go holdConnection(conn)
	}
}

func AppendIfMissing(slice []message.Node, i message.Node) []message.Node {
	for _, ele := range slice {
		if ele.IP == i.IP {
			return slice
		}
	}
	newNodes <- i
	newNodesPaxos <- i
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

func holdConnection(conn net.Conn) {
	var connectionOK = true
	for connectionOK == true {
		decoder := gob.NewDecoder(conn)
		var msg interface{}
		var err = decoder.Decode(&msg)
		if err != nil {
			connectionOK = false
		}
		if msg != nil {
			var clientMsg message.ClientRequestMessage
			clientMsg = msg.(message.ClientRequestMessage)
			if leader.IP == selfnode.IP {
				if clientMsg.RemoteAddress == "" {
					RoundVar.GetRound().RespondClient = Utils.GetIp(conn.RemoteAddr().String())
				} else {
					RoundVar.GetRound().RespondClient = clientMsg.RemoteAddress
				}
				acceptorChan <- clientMsg.Content
			} else {
				fmt.Println("Im not leader, sending it on!")
				if clientMsg.RemoteAddress == "" {
					clientMsg.RemoteAddress = Utils.GetIp(conn.RemoteAddr().String())
				}
				leaderService := leader.IP + ":1337"
				fmt.Println(leaderService)
				leaderCon, err := net.Dial("tcp", leaderService)
				if err != nil {
					fmt.Println(err)
				} else {
					encoder := gob.NewEncoder(leaderCon)
					msg = clientMsg
					encoder.Encode(&msg)
				}
				leaderCon.Close()
			}
		}
	}
	fmt.Println("Client closed connection, no more to share")
}
