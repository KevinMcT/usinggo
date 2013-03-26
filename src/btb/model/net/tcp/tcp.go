package tcp

import (
	"btb/controller/node"
	"btb/model/net/msg"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

var (
	restart bool
)

func Listen(nodeChan chan []node.T_Node, tcpLeaderRequestChan chan node.T_Node, machineCountChan chan msg.MACHINECOUNT, messageChan chan string, tcpLeaderResponseChan chan node.T_Node, tcpHartBeatRequest chan msg.HARTBEATREQUEST, tcpHartBeatResponse chan msg.HARTBEATRESPONSE, leaderDown chan int) {
	service := "0.0.0.0:2000"
	tcpAddr, _ := net.ResolveTCPAddr("tcp", service)
	var message interface{}
	var nodeList []node.T_Node
	var leaderResponse node.T_Node
	var leaderRequest node.T_Node
	for {
		listener, err := net.ListenTCP("tcp", tcpAddr)
		//listener.SetDeadline(time.Now().Add(10 * time.Second))
		err = nil
		var inMessage string
		var machineCount msg.MACHINECOUNT
		var conn *net.TCPConn
		for {
			conn, err = listener.AcceptTCP()
			if err != nil {
				fmt.Println("Error in TCP", err)
				conn.Close()
				listener.Close()
				err = nil
				fmt.Println("Break")
				break
			}
			decoder := gob.NewDecoder(conn)
			decoder.Decode(&message)
			switch message.(type) {
			case msg.LISTRESPONSE:
				nodeList = message.(msg.LISTRESPONSE).LIST
				//listener.SetDeadline(time.Now().Add(5 * time.Second))
				nodeChan <- nodeList
			case msg.HARTBEATREQUEST:
				tcpHartBeatRequest <- message.(msg.HARTBEATREQUEST)
				listener.SetDeadline(time.Now().Add(250 * time.Millisecond))
			case msg.HARTBEATRESPONSE:
				tcpHartBeatResponse <- message.(msg.HARTBEATRESPONSE)
				//listener.SetDeadline(time.Now().Add(60 * time.Second))
			case msg.LEADERREQUEST:
				leaderRequest = message.(msg.LEADERREQUEST).FROMNODE
				//listener.SetDeadline(time.Now().Add(250 * time.Millisecond))
				tcpLeaderRequestChan <- leaderRequest
			case msg.LEADERRESPONSE:
				leaderResponse = message.(msg.LEADERRESPONSE).NODE
				//listener.SetDeadline(time.Now().Add(1000 * time.Millisecond))
				tcpLeaderResponseChan <- leaderResponse
			case msg.Node:
			case msg.Lead:
			case msg.MACHINECOUNT:
				machineCount = message.(msg.MACHINECOUNT)
				//listener.SetDeadline(time.Now().Add(1000 * time.Millisecond))
				machineCountChan <- machineCount
			case msg.MESSAGE:
				inMessage = message.(msg.MESSAGE).MSG
				messageChan <- inMessage
			}
		}
		//listener.SetDeadline(time.Now().Add(10 * time.Second))
		leaderDown <- 1
		fmt.Println("Sent leader down channel")
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
