package core

import (
	"encoding/gob"
	"fmt"
	"lab6/Utils"
	"lab6/controller/node"
	"lab6/model/RoundVar"
	"lab6/model/SlotList"
	"lab6/model/fd"
	"lab6/model/le"
	"lab6/model/net/msg"
	"lab6/model/net/tcp"
	"lab6/model/net/udp"
	"lab6/model/px"
	"net"
	"os"
	"time"
)

var (
	udpListenChan           = make(chan string, 10)
	createNodeChan          = make(chan node.T_Node, 10)
	nodeList                = make([]node.T_Node, 0)
	me                      node.T_Node
	leader                  node.T_Node
	tcpNodeChan             = make(chan []node.T_Node, 10)
	tcpLeaderRequestChan    = make(chan node.T_Node, 10)
	tcpHartbeatRequestChan  = make(chan msg.HARTBEATREQUEST, 10)
	tcpHartbeatResponseChan = make(chan msg.HARTBEATRESPONSE, 10)
	machineCountChan        = make(chan msg.MACHINECOUNT, 10)
	messageChan             = make(chan string, 10)
	leaderDown              = make(chan int, 1)
	tcpLeaderResponseChan   = make(chan node.T_Node, 10)

	newNodeChan     = make(chan node.T_Node, 10)
	newLeaderChan   = make(chan node.T_Node, 1)
	suspectedChan   = make(chan node.T_Node, 10)
	restoreChan     = make(chan node.T_Node, 10)
	endNodeListChan = make(chan []node.T_Node, 10)

	newNodesPaxos = make(chan node.T_Node, 10)
	acceptorChan  = make(chan string, 20)
	slots         = SlotList.NewSlots()
)

func Paxos() {
	fmt.Println("Starting software. . .")
	var startTime = time.Now().Unix()
	fmt.Println("Starting node creation. . . ")
	go node.Node(udpListenChan, createNodeChan)
	fmt.Println("--Done!--")
	fmt.Println("Starting TCP listen. . .")
	go tcp.Listen(tcpNodeChan, tcpLeaderRequestChan, machineCountChan, messageChan, tcpLeaderResponseChan, tcpHartbeatRequestChan, tcpHartbeatResponseChan, leaderDown)
	fmt.Println("--Done!--")
	fmt.Println("Starting dynamic nodeadding. . .")
	go addNodesFromUdp(udpListenChan)
	fmt.Println("--Done!--")
	fmt.Println("Starting UDP listen. . .")
	go udp.Listen(udpListenChan)
	fmt.Println("--Done!--")
	me = createMe(startTime)
	nodeList = append(nodeList, me)
	fmt.Println("Sending UDP broadcast. . .")
	udp.SendBroadcast(startTime)
	fmt.Println("--Done!--")
	if me.LEAD == false {
		fmt.Println("Waiting for node list to elect first leader. . .")
		nodeList = <-tcpNodeChan
		fmt.Println("--Done!--")
	}
	fmt.Println("Electing first leader. . .")
	leader = le.Elect(nodeList)
	RoundVar.GetRound().CurrentLeader = leader
	for i, v := range nodeList {
		if v.IP == leader.IP {
			nodeList[i].LEAD = true
			fmt.Println("--Done!--")
		}
	}
	fmt.Println("--Elected ", leader, "as leader--")
	fmt.Println("---------------------------------")
	fmt.Println("--Starting Failuredetect for all nodes. . . ")
	go fd.Detect(me, leader, newLeaderChan, newNodeChan, suspectedChan, restoreChan, tcpHartbeatRequestChan, tcpHartbeatResponseChan, nodeList, endNodeListChan)
	fmt.Println("--Done!--")

	//Paxos go routines
	go ClientConnection()
	go px.Acceptor()
	go px.Learner()
	go px.PaxosHandler()
	go px.Proposer(me, newNodesPaxos, acceptorChan)

	//Main program loop
	for {
		var ever = true
		for ever {
			timeout := make(chan bool, 1)
			go func() {
				time.Sleep(5000 * time.Millisecond)
				timeout <- true
			}()
			select {
			case suspectedNode := <-suspectedChan:
				fmt.Println("--Suspect on node: ", suspectedNode.IP, "--")
				for i, v := range nodeList {
					if v.IP == suspectedNode.IP {
						nodeList[i].SUSPECTED = true
						newNodeChan <- nodeList[i]
					}
					if me.IP != v.IP {
						go tcp.Send(v.IP, msg.LISTRESPONSE{LIST: nodeList})
					}
					endNodeListChan <- nodeList
				}
			case restoredNode := <-restoreChan:
				fmt.Println("--Restore on node: ", restoredNode.IP, "--")
				for i, v := range nodeList {
					if v.IP == restoredNode.IP {
						nodeList[i].SUSPECTED = false
						newNodeChan <- nodeList[i]
					}
					if me.IP != v.IP {
						go tcp.Send(v.IP, msg.LISTRESPONSE{LIST: nodeList})
					}
					endNodeListChan <- nodeList
				}
			case <-leaderDown:
				fmt.Println("--Suspected leader down--")
				for i, v := range nodeList {
					if v.IP == leader.IP {
						nodeList[i].LEAD = false
						nodeList[i].SUSPECTED = true
						newNodeChan <- nodeList[i]
					}
				}
				leader = le.Elect(nodeList)
				for i, v := range nodeList {
					if v.IP == leader.IP {
						nodeList[i].LEAD = true
						newLeaderChan <- leader
						fmt.Println("--Done!--")
					}
				}
				endNodeListChan <- nodeList
				fmt.Println("--Elected ", leader, "as leader--")
				RoundVar.GetRound().CurrentLeader = leader
				//TODO Need to make this unique for each node, so node
				//node can propose the same round.
				RoundVar.GetRound().Round = RoundVar.GetRound().Round + 1
				msg.RestartProposer <- "newLeader"
				if me.IP == leader.IP {
					go udp.Listen(udpListenChan)
				}
			case list := <-tcpNodeChan:
				nodeList = list
				RoundVar.GetRound().List = nodeList
				endNodeListChan <- list
			case <-timeout:
				for _, v := range nodeList {
					if me.IP != v.IP {
						go tcp.Send(v.IP, msg.LISTRESPONSE{LIST: nodeList})
					}
				}
			}
		}

	}
}

func addNodesFromUdp(inputChan chan string) {
	var exists bool
	exists = false
	for {
		var node node.T_Node
		node = <-createNodeChan
		me.LEAD = true
		for i, v := range nodeList {
			if me.IP == v.IP {
				nodeList[i].LEAD = true
			}
		}
		for i, v := range nodeList {
			if v.IP == node.IP {
				exists = true
				fmt.Println("--Node already in system--")
				restoreChan <- node
				nodeList[i].SUSPECTED = false
				nodeList[i].TIME = node.TIME
				break
			}
		}
		if !exists {
			if node.IP != me.IP {
				fmt.Println("--New node detected on address:", node.IP, "--")
			}
			newNodeChan <- node
			newNodesPaxos <- node
			nodeList = append(nodeList, node)
			RoundVar.GetRound().List = nodeList
			if len(nodeList) > 2 && leader.IP == me.IP {
				msg.SendPrepareChan <- true
			}
		}
		exists = false
		if node.IP != me.IP {
			for _, v := range nodeList {
				go tcp.Send(v.IP, msg.LISTRESPONSE{LIST: nodeList})
			}
		}
	}
}

func createMe(startTime int64) node.T_Node {
	var selfnode node.T_Node
	name, _ := os.Hostname()
	addr, _ := net.LookupHost(name)
	UDPAddr, _ := net.ResolveUDPAddr("udp4", addr[0]+":1888")
	selfnode.IP = UDPAddr.IP.String()
	selfnode.TIME = startTime
	return selfnode
}

/*
Method for listning after client requests. If this instance is the
leader of the paxos system he sends the message on to the proposer. If
not the message is redirected to the node that is the leader.
*/
func ClientConnection() {
	fmt.Println("Waiting for inncoming clients. . . ")
	service := "0.0.0.0:1337"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	Utils.CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	for {
		Utils.CheckError(err)
		conn, _ := listener.Accept()
		go holdReplicaConnection(conn)
	}
}

func contains(s []msg.Node, e bool) bool {
	for _, a := range s {
		if a.LEAD == e {
			return true
		}
	}
	return false
}

func holdReplicaConnection(conn net.Conn) {
	var connectionOK = true
	for connectionOK == true {
		decoder := gob.NewDecoder(conn)
		var message interface{}
		var err = decoder.Decode(&message)
		if err != nil {
			connectionOK = false
			fmt.Println("--Error in Paxos: ", err, "--")
		}
		if message != nil {
			switch message.(type) {
			case msg.ClientRequestMessage:
				var clientMsg msg.ClientRequestMessage
				clientMsg = message.(msg.ClientRequestMessage)
				if leader.IP == me.IP {
					if clientMsg.RemoteAddress == "" {
						RoundVar.GetRound().RespondClient = Utils.GetIp(conn.RemoteAddr().String())
					} else {
						RoundVar.GetRound().RespondClient = clientMsg.RemoteAddress
					}
					acceptorChan <- clientMsg.Content
				} else {
					fmt.Println("--Not leader. Relay message to: ", leader.IP, " --")
					if clientMsg.RemoteAddress == "" {
						clientMsg.RemoteAddress = Utils.GetIp(conn.RemoteAddr().String())
					}
					leaderService := leader.IP + ":1337"
					leaderCon, err := net.Dial("tcp", leaderService)
					if err != nil {
						fmt.Println("--Paxos:", err, "--")
					} else {
						encoder := gob.NewEncoder(leaderCon)
						message = clientMsg
						encoder.Encode(&message)
					}
				}
			case msg.ClientRequestNodes:
				var clientMsg msg.ClientRequestNodes
				clientMsg = message.(msg.ClientRequestNodes)
				if leader.IP == me.IP {
					tcp.SendPaxosMessage(clientMsg.RemoteAddress, msg.ClientResponseNodes{List: nodeList})
				} else {
					fmt.Println("--Not leader. Relay message to: ", leader.IP, " --")
					if clientMsg.RemoteAddress == "" {
						clientMsg.RemoteAddress = Utils.GetIp(conn.RemoteAddr().String())
					}
					leaderService := leader.IP + ":1337"
					leaderCon, err := net.Dial("tcp", leaderService)
					if err != nil {
						fmt.Println("--Paxos:", err, "--")
					} else {
						encoder := gob.NewEncoder(leaderCon)
						message = clientMsg
						encoder.Encode(&message)
					}
				}
			}
		} else {
			fmt.Println("--Sending empty messages!--")
		}
	}
	fmt.Println("--Client closed connection, no more to share--")
}
