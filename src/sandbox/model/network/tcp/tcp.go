package tcp

import (
	"encoding/gob"
	"fmt"
	"net"
	"sandbox/controller/node"
	"sandbox/model/network/message"
	"time"
)

var (
	restart bool
)

func Listen(nodeChan chan []node.T_Node, tcpLeaderRequestChan chan node.T_Node, machineCountChan chan message.MACHINECOUNT, messageChan chan string, tcpLeaderResponseChan chan node.T_Node, tcpHartBeatRequest chan message.HARTBEATREQUEST, tcpHartBeatResponse chan message.HARTBEATRESPONSE, leaderDown chan int) {
	service := "0.0.0.0:2000"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	var msg interface{}
	var nodeList []node.T_Node
	var leaderResponse node.T_Node
	var leaderRequest node.T_Node
	for !restart {
		restart = false
		var inMessage string
		var machineCount message.MACHINECOUNT
		var conn *net.TCPConn
		for {
			conn, err = listener.AcceptTCP()
			if err != nil {
				fmt.Println("WHOAH")
				conn.Close()
				listener.Close()
				leaderDown <- 1
				err = nil
				restart = true
				break
			}
			decoder := gob.NewDecoder(conn)
			decoder.Decode(&msg)
			switch msg.(type) {
			case message.LISTRESPONSE:
				nodeList = msg.(message.LISTRESPONSE).LIST
				nodeChan <- nodeList
			case message.HARTBEATREQUEST:
				tcpHartBeatRequest <- msg.(message.HARTBEATREQUEST)
				listener.SetDeadline(time.Now().Add(700 * time.Millisecond))
			case message.HARTBEATRESPONSE:
				tcpHartBeatResponse <- msg.(message.HARTBEATRESPONSE)
			case message.LEADERREQUEST:
				leaderRequest = msg.(message.LEADERREQUEST).FROMNODE
				//fmt.Println("Request from node:", leaderRequest)
				tcpLeaderRequestChan <- leaderRequest
			case message.LEADERRESPONSE:
				leaderResponse = msg.(message.LEADERRESPONSE).NODE
				//fmt.Println("Response from leader:", leaderResponse)
				tcpLeaderResponseChan <- leaderResponse
			case message.Node:
			case message.Lead:
			case message.MACHINECOUNT:
				machineCount = msg.(message.MACHINECOUNT)
				fmt.Println("MachineCount: ", machineCount)
				machineCountChan <- machineCount
			case message.MESSAGE:
				inMessage = msg.(message.MESSAGE).MSG
				messageChan <- inMessage
			}
			conn.Close()
		}
	}
}

func Send(ip string, msg interface{}) error {
	service := ip + ":2000"
	conn, err := net.Dial("tcp", service)
	if err != nil {
		return err
	} else {
		encoder := gob.NewEncoder(conn)
		encoder.Encode(&msg)
	}
	conn.Close()
	return err
}

func AppendIfMissing(slice []node.T_Node, i node.T_Node) []node.T_Node {
	for _, ele := range slice {
		if ele.IP == i.IP {
			return slice
		}
	}
	return append(slice, i)
}
