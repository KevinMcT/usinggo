package main

import (
	"fmt"
	"net"
	"os"
	"sandbox/controller/machine"
	"sandbox/controller/node"
	"sandbox/model/failuredetect"
	"sandbox/model/leaderelection"
	"sandbox/model/network/message"
	"sandbox/model/network/tcp"
	"sandbox/model/network/udp"
	"time"
)

var (
	udpListenChan           = make(chan string, 10)
	udpMasterChan           = make(chan string, 10)
	machineChan             = make(chan machine.T_Machine, 10)
	leaderElectionChan      = make(chan []node.T_Node, 10)
	leaderElectedChan       = make(chan node.T_Node, 10)
	leaderBlock             = make(chan int, 10)
	tcpNodeChan             = make(chan []node.T_Node, 10)
	tcpLeaderRequestChan    = make(chan node.T_Node, 10)
	tcpHartbeatRequestChan  = make(chan message.HARTBEATREQUEST, 10)
	tcpHartbeatResponseChan = make(chan message.HARTBEATRESPONSE, 10)
	suspectedChan           = make(chan node.T_Node, 10)
	restoreChan             = make(chan node.T_Node, 10)
	machineCountChan        = make(chan message.MACHINECOUNT, 10)
	messageChan             = make(chan string, 10)
	gotAllMachinesChan      = make(chan node.T_Node, 10)
	tcpLeaderResponseChan   = make(chan node.T_Node, 10)
	newNodeChan             = make(chan node.T_Node, 10)
	leaderDown              = make(chan int, 1)
	newLeaderChan           = make(chan node.T_Node, 1)
	askForLeaderChan        = make(chan int, 1)
	endNodeListChan         = make(chan []node.T_Node, 10)
	machineList             = make([]machine.T_Machine, 0)
	machineCountList        = make([]message.MACHINECOUNT, 0)
	nodeList                = make([]node.T_Node, 0)
	tcpLeaderList           = make([]node.T_Node, 0)
	listChanged             []node.T_Node
	leader                  node.T_Node
	machineCount            message.MACHINECOUNT
	inMessage               string
	me                      machine.T_Machine
)

func main() {
	go tcp.Listen(tcpNodeChan, tcpLeaderRequestChan, machineCountChan, messageChan, tcpLeaderResponseChan, tcpHartbeatRequestChan, tcpHartbeatResponseChan, leaderDown)
	fmt.Println("Go")
	startTime := time.Now().UnixNano()
	fmt.Println("Start Listen")
	go udp.Listen(udpListenChan, udpMasterChan)
	go machine.Machine(udpListenChan, machineChan)
	fmt.Println("Send broadcast")
	udp.SendBroadcast(startTime)
	me = CreateSelf(startTime)
	machineList, _ = AppendIfMissing(machineList, me)
	go fillList()
	go askForLeader()
	go leaderResponse()
	for {
		if leader.IP != "" && leader.IP == me.IP {
			fmt.Println("Started failure for leader")
			go failuredetect.Detect(me, leader, newLeaderChan, newNodeChan, suspectedChan, restoreChan, tcpHartbeatRequestChan, tcpHartbeatResponseChan, nodeList, endNodeListChan)
			break
		}
		if leader.IP != "" && leader.IP != me.IP {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	for {
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(10 * time.Millisecond)
			timeout <- true
		}()
		select {
		case suspectedNode := <-suspectedChan:
			for i, v := range nodeList {
				if v.IP == suspectedNode.IP {
					fmt.Println("Suspect on node:", suspectedNode.IP)
					fmt.Println("TIME: ", suspectedNode.TIME)
					nodeList[i].SUSPECTED = true
				}
				if v.IP != me.IP && v.SUSPECTED != true {
					fmt.Println("SENDING LIST AFTER SUSPECT")
					tcp.Send(v.IP, message.LISTRESPONSE{LIST: nodeList})
				}
			}
		case restoreNode := <-restoreChan:
			for i, v := range nodeList {
				if v.IP == restoreNode.IP {
					fmt.Println("Restore on node:", restoreNode.IP)
					fmt.Println("TIME:", restoreNode.TIME)
					nodeList[i].SUSPECTED = false
					nodeList[i].LEAD = false
					newNodeChan <- v
					tcp.Send(v.IP, message.MACHINECOUNT{I: len(nodeList), NODE: nodeList})
				}
				if v.IP != me.IP {
					tcp.Send(v.IP, message.LISTRESPONSE{LIST: nodeList})
				}
			}
		case <-leaderDown:
			fmt.Println("Suspected leader down")
			fmt.Println("----------------------------------------------------------------")
			for i, v := range nodeList {
				if v.IP == leader.IP {
					nodeList[i].SUSPECTED = true
					nodeList[i].LEAD = false
				}
			}
			fmt.Println("Nodelist after leader down: ", nodeList)
			leader = leaderelection.Elect(nodeList)
			for i, v := range nodeList {
				if me.IP == leader.IP && me.IP == v.IP {
					nodeList[i].LEAD = true
				}
			}
			fmt.Println("Leader elect in main after break: ", leader)
			//time.Sleep(5 * time.Millisecond)
			for {
				if leader.IP != "" && leader.IP == me.IP {
					fmt.Println("Hello leader <3") //My leader is Tina N.D, but for now this node will do....
					udpListenChan = make(chan string, 10)
					udpMasterChan = make(chan string, 10)
					go udp.Listen(udpListenChan, udpMasterChan)
					fmt.Println("UDP operational!")
					newLeaderChan <- leader
					endNodeListChan <- nodeList
					fmt.Println("After new leader is sent to FD")
					break
				}
				if leader.IP != "" && leader.IP != me.IP {
					fmt.Println("I should ask for leader here...")
					newLeaderChan <- leader
					endNodeListChan <- nodeList
					askForLeaderChan <- 1
					break
				}
				//time.Sleep(2 * time.Millisecond)
			}

		case <-timeout:

		case listChanged = <-tcpNodeChan:
			if leader.IP != me.IP {
				nodeList = listChanged
				//fmt.Println("Recieved list:", nodeList)
			}
		}
	}

}

func leaderResponse() {
	for {
		//fmt.Println("Leader response")
		leaderRequestMsg := <-tcpLeaderRequestChan
		if leader.IP == "" {
			leader = leaderelection.Elect(nodeList)
		}
		//fmt.Println("Sending leader response", leader)
		go tcp.Send(leaderRequestMsg.IP, message.LEADERRESPONSE{NODE: leader})
		//time.Sleep(10 * time.Millisecond)
		//fmt.Println("SENDING LIST TO REQUESTER:", leaderRequestMsg.IP)
		go tcp.Send(leaderRequestMsg.IP, message.LISTRESPONSE{LIST: nodeList})
		me.LEAD = true
		for i, v := range nodeList {
			if me.IP == v.IP {
				nodeList[i].LEAD = true
			}
		}
	}
}

func fillList() {
	var i bool
	i = false
	for {
		var mac machine.T_Machine
		mac = <-machineChan
		machineList, i = AppendIfMissing(machineList, mac)
		if i {
			for _, v := range nodeList {
				fmt.Println(v)
				if v.IP != me.IP {
					fmt.Println("SENDING LIST IN fillLIST")
					go tcp.Send(v.IP, message.LISTRESPONSE{LIST: nodeList})
				}
			}
		}

		go tcp.Send(mac.IP, message.MACHINECOUNT{I: len(nodeList), NODE: nodeList})
		newNodeChan <- node.CreateNode(mac)
	}
}

func askForLeader() {
	for {
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(50 * time.Millisecond)
			timeout <- true
		}()
		select {
		case <-askForLeaderChan:
			fmt.Println("TEST")
			var lead node.T_Node
			if leader.IP != me.IP {
				go tcp.Send(leader.IP, message.LEADERREQUEST{TONODE: leader, FROMNODE: node.CreateNode(me)})
				lead = <-tcpLeaderResponseChan
				//newLeaderChan <- lead
				fmt.Println("Confirmed leader: ", lead.IP)
			}
		case <-timeout:
			var mac message.MACHINECOUNT
			var lead node.T_Node
			timeout1 := make(chan bool, 1)
			go func() {
				time.Sleep(5 * time.Millisecond)
				timeout1 <- true
			}()
			select {
			case mac = <-machineCountChan:
				if leader.IP != me.IP {
					leader = leaderelection.Elect(mac.NODE)
					fmt.Println("Leader: ", leader)
					go tcp.Send(leader.IP, message.LEADERREQUEST{TONODE: leader, FROMNODE: node.CreateNode(me)})
					lead = <-tcpLeaderResponseChan
					fmt.Println("Confirmed leader: ", lead.IP)
					go failuredetect.Detect(me, leader, newLeaderChan, newNodeChan, suspectedChan, restoreChan, tcpHartbeatRequestChan, tcpHartbeatResponseChan, nodeList, endNodeListChan)
				}
			case <-timeout1:
			}
		}
	}
}

func AppendIfMissing(slice []machine.T_Machine, i machine.T_Machine) ([]machine.T_Machine, bool) {
	for _, ele := range slice {
		if ele.IP == i.IP {
			fmt.Println("Node already in system")
			for j, v := range nodeList {
				if i.IP == v.IP {
					fmt.Println("Old List:", nodeList[j])
					nodeList[j].SUSPECTED = false
					nodeList[j].TIME = i.TIME
					fmt.Println("New List:", nodeList[j])
					restoreChan <- nodeList[j]
				}
			}
			return slice, false
		}
	}
	node := node.CreateNode(i)
	if node.IP != me.IP {
		fmt.Println("New node detected on address:", node.IP)
	}
	nodeList = append(nodeList, node)
	return append(slice, i), true
}

func CreateSelf(startTime int64) machine.T_Machine {
	var selfnode machine.T_Machine
	name, _ := os.Hostname()
	addr, _ := net.LookupHost(name)
	UDPAddr, _ := net.ResolveUDPAddr("udp4", addr[0]+":1888")
	selfnode.IP = UDPAddr.IP.String()
	selfnode.TIME = startTime
	return selfnode
}
