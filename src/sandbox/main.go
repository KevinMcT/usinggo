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
	leaderElectionChan      = make(chan []node.T_Node, 0)
	leaderElectedChan       = make(chan node.T_Node, 0)
	leaderBlock             = make(chan int, 0)
	tcpNodeChan             = make(chan []node.T_Node, 10)
	tcpLeaderRequestChan    = make(chan node.T_Node, 0)
	tcpHartbeatRequestChan  = make(chan message.HARTBEATREQUEST, 0)
	tcpHartbeatResponseChan = make(chan message.HARTBEATRESPONSE, 0)
	suspectedChan           = make(chan node.T_Node, 0)
	restoreChan             = make(chan node.T_Node, 0)
	machineCountChan        = make(chan message.MACHINECOUNT, 0)
	messageChan             = make(chan string, 0)
	gotAllMachinesChan      = make(chan node.T_Node, 0)
	tcpLeaderResponseChan   = make(chan node.T_Node, 0)
	newNodeChan             = make(chan node.T_Node, 0)
	leaderDown              = make(chan int, 1)
	newLeaderChan           = make(chan node.T_Node, 1)
	askForLeaderChan        = make(chan int, 1)
	machineList             = make([]machine.T_Machine, 0)
	machineCountList        = make([]message.MACHINECOUNT, 0)
	nodeList                = make([]node.T_Node, 0)
	tcpLeaderList           = make([]node.T_Node, 0)
	leader                  node.T_Node
	machineCount            message.MACHINECOUNT
	inMessage               string
	me                      machine.T_Machine
)

func main() {
	fmt.Println("Go")
	startTime := time.Now().UnixNano()
	fmt.Println("Start Listen")
	go udp.Listen(udpListenChan, udpMasterChan)
	go machine.Machine(udpListenChan, machineChan)
	go tcp.Listen(tcpNodeChan, tcpLeaderRequestChan, machineCountChan, messageChan, tcpLeaderResponseChan, tcpHartbeatRequestChan, tcpHartbeatResponseChan, leaderDown)
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
			go failuredetect.Detect(me, leader, newLeaderChan, newNodeChan, suspectedChan, restoreChan, tcpHartbeatRequestChan, tcpHartbeatResponseChan, nodeList)
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
					nodeList[i].SUSPECTED = true
				}
			}
		case restoreNode := <-restoreChan:
			for i, v := range nodeList {
				if v.IP == restoreNode.IP {
					fmt.Println("Restore on node:", restoreNode.IP)
					nodeList[i].SUSPECTED = false
				}
			}
		case <-leaderDown:
			fmt.Println("Suspected leader down")
			for i, v := range nodeList {
				if v.IP == leader.IP {
					nodeList[i].SUSPECTED = true
					nodeList[i].LEAD = false
				}
			}
			fmt.Println("Nodelist after leader down: ", nodeList)
			leader = leaderelection.Elect(nodeList, leaderElectedChan, leaderBlock)
			for i, v := range nodeList {
				if v.LEAD == true && v.SUSPECTED == true {
					fmt.Println("Safety!")
					nodeList[i].LEAD = false
				}
			}
			fmt.Println("Leader elect in main after break: ", leader)
			time.Sleep(5 * time.Millisecond)
			for {
				if leader.IP != "" && leader.IP == me.IP {
					fmt.Println("Hello leader <3") //My leader is Tina N.D, but for now this node will do....
					udpListenChan = make(chan string, 10)
					udpMasterChan = make(chan string, 10)
					go udp.Listen(udpListenChan, udpMasterChan)
					fmt.Println("UDP operational!")
					newLeaderChan <- leader
					for _, v := range nodeList {
						go tcp.Send(v.IP, message.MACHINECOUNT{I: len(nodeList), NODE: nodeList})
					}
					fmt.Println("After new leader is sent to FD")
					break
				}
				if leader.IP != "" && leader.IP != me.IP {
					fmt.Println("I should ask for leader here...")
					askForLeaderChan <- 1
					break
				}
				time.Sleep(5 * time.Millisecond)
			}

		case <-timeout:

		case listChanged := <-tcpNodeChan:
			fmt.Println("list changed")
			nodeList = listChanged
			fmt.Println(nodeList)
		}

		time.Sleep(1 * time.Millisecond)
	}

}

func leaderResponse() {
	for {
		fmt.Println("Leader response")
		leaderRequestMsg := <-tcpLeaderRequestChan
		leader = leaderelection.Elect(nodeList, leaderElectedChan, leaderBlock)
		fmt.Println("Sending leader response")
		tcp.Send(leaderRequestMsg.IP, message.LEADERRESPONSE{NODE: leader})
		time.Sleep(10 * time.Millisecond)
		tcp.Send(leaderRequestMsg.IP, message.LISTRESPONSE{LIST: nodeList})
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
					err := tcp.Send(v.IP, message.LISTRESPONSE{LIST: nodeList})
					fmt.Println(err)
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
			time.Sleep(5 * time.Millisecond)
			timeout <- true
		}()
		select {
		case <-askForLeaderChan:
			fmt.Println("TEST")
			var lead node.T_Node
			leader = leaderelection.Elect(nodeList, leaderElectedChan, leaderBlock)
			fmt.Println("New leader after leader crach: ", leader)
			newLeaderChan <- leader
			tcp.Send(leader.IP, message.LEADERREQUEST{TONODE: leader, FROMNODE: node.CreateNode(me)})
			lead = <-tcpLeaderResponseChan
			fmt.Println("Confirmed leader: ", lead.IP)
			fmt.Println("Leaders list:", nodeList)
		case <-timeout:
			var mac message.MACHINECOUNT
			var lead node.T_Node
			timeout = make(chan bool, 1)
			go func() {
				time.Sleep(5 * time.Millisecond)
				timeout <- true
			}()
			select {
			case mac = <-machineCountChan:
				leader = leaderelection.Elect(mac.NODE, leaderElectedChan, leaderBlock)
				fmt.Println("Leader: ", leader)
				tcp.Send(leader.IP, message.LEADERREQUEST{TONODE: leader, FROMNODE: node.CreateNode(me)})
				lead = <-tcpLeaderResponseChan
				fmt.Println("Confirmed leader: ", lead.IP)
				fmt.Println("Leaders list:", nodeList)
				go failuredetect.Detect(me, leader, newLeaderChan, newNodeChan, suspectedChan, restoreChan, tcpHartbeatRequestChan, tcpHartbeatResponseChan, nodeList)
			case <-timeout:
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
					nodeList[j].SUSPECTED = false
					nodeList[j].TIME = i.TIME
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
