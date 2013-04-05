package core

import (
	"encoding/gob"
	"fmt"
	"lab5Merge/Utils"
	"lab5Merge/controller/node"
	"lab5Merge/model/RoundVar"
	"lab5Merge/model/fd"
	"lab5Merge/model/le"
	"lab5Merge/model/net/msg"
	"lab5Merge/model/net/tcp"
	"lab5Merge/model/net/udp"
	"lab5Merge/model/px"
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
)

func Paxos() {
	fmt.Println("Starting software. . .")
	var startTime = time.Now().Unix()
	fmt.Println("Starting node creation. . . ")
	go node.Node(udpListenChan, createNodeChan)
	fmt.Println("--Done!--")
	fmt.Println("Starting dynamic nodeadding. . .")
	go addNodesFromUdp(udpListenChan)
	fmt.Println("--Done!--")
	fmt.Println("Starting TCP listen. . .")
	go tcp.Listen(tcpNodeChan, tcpLeaderRequestChan, machineCountChan, messageChan, tcpLeaderResponseChan, tcpHartbeatRequestChan, tcpHartbeatResponseChan, leaderDown)
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
	for i, v := range nodeList {
		if v.IP == leader.IP {
			nodeList[i].LEAD = true
			fmt.Println("--Done!--")
		}
	}
	fmt.Println("--Elected ", leader, "as leader--")
	fmt.Println("---------------------------------")
	fmt.Println("Starting Failuredetect for all nodes. . . ")
	go fd.Detect(me, leader, newLeaderChan, newNodeChan, suspectedChan, restoreChan, tcpHartbeatRequestChan, tcpHartbeatResponseChan, nodeList, endNodeListChan)
	fmt.Println("--Done!--")

	//Paxos go routines
	go ClientConnection()
	go px.Acceptor()
	go px.Learner()
	go px.PaxosHandler()
	go px.Proposer(leader, me, newNodesPaxos, acceptorChan)

	//Main program loop
	for {
		var ever = true
		for ever {
			timeout := make(chan bool, 1)
			go func() {
				time.Sleep(1 * time.Millisecond)
				timeout <- true
			}()
			select {
			case suspectedNode := <-suspectedChan:
				fmt.Println("Suspect on node: ", suspectedNode.IP)
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
				fmt.Println("Restore on node: ", restoredNode.IP)
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
				if me.IP == leader.IP {
					go udp.Listen(udpListenChan)
				}
			case list := <-tcpNodeChan:
				nodeList = list
				endNodeListChan <- list
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
				fmt.Println("Going to send prepare!!")
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
	fmt.Println("Waiting for inncoming clients")
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
			fmt.Println("Error in Paxos: ", err)
		}
		if message != nil {
			var clientMsg msg.ClientRequestMessage
			clientMsg = message.(msg.ClientRequestMessage)
			fmt.Println("ClientMessage: ", clientMsg)
			if leader.IP == me.IP {
				if clientMsg.RemoteAddress == "" {
					RoundVar.GetRound().RespondClient = Utils.GetIp(conn.RemoteAddr().String())
				} else {
					RoundVar.GetRound().RespondClient = clientMsg.RemoteAddress
				}
				fmt.Println("--Sending message on acceptor chan--")
				acceptorChan <- clientMsg.Content
				fmt.Println("--Sent message on acceptor chan--")
			} else {
				fmt.Println("Im not leader, sending it on!")
				if clientMsg.RemoteAddress == "" {
					clientMsg.RemoteAddress = Utils.GetIp(conn.RemoteAddr().String())
				}
				leaderService := leader.IP + ":1337"
				fmt.Println(leaderService)
				leaderCon, err := net.Dial("tcp", leaderService)
				if err != nil {
					fmt.Println("Paxos:", err)
				} else {
					encoder := gob.NewEncoder(leaderCon)
					message = clientMsg
					encoder.Encode(&message)
				}
			}
		} else {
			fmt.Println("Sending empty messages stupid!")
		}
	}
	//tcp.StoreDecoder(conn, *decoder)
	fmt.Println("Client closed connection, no more to share")
}
