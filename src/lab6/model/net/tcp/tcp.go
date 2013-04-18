package tcp

import (
	"encoding/gob"
	"fmt"
	"lab6/controller/node"
	"lab6/model/net/msg"
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
		err = nil
		var inMessage string
		var machineCount msg.MACHINECOUNT
		var conn *net.TCPConn
		for {
			conn, err = listener.AcceptTCP()
			if err != nil {
				fmt.Println("--Error in TCP", err, "--")
				conn.Close()
				listener.Close()
				err = nil
				break
			}
			decoder := gob.NewDecoder(conn)
			decoder.Decode(&message)
			switch message.(type) {
			case msg.LISTRESPONSE:
				nodeList = message.(msg.LISTRESPONSE).LIST
				nodeChan <- nodeList
			case msg.HARTBEATREQUEST:
				tcpHartBeatRequest <- message.(msg.HARTBEATREQUEST)
				listener.SetDeadline(time.Now().Add(500 * time.Millisecond))
			case msg.HARTBEATRESPONSE:
				tcpHartBeatResponse <- message.(msg.HARTBEATRESPONSE)
			case msg.LEADERREQUEST:
				leaderRequest = message.(msg.LEADERREQUEST).FROMNODE
				tcpLeaderRequestChan <- leaderRequest
			case msg.LEADERRESPONSE:
				leaderResponse = message.(msg.LEADERRESPONSE).NODE
				tcpLeaderResponseChan <- leaderResponse
			case msg.Node:
			case msg.Lead:
			case msg.MACHINECOUNT:
				machineCount = message.(msg.MACHINECOUNT)
				machineCountChan <- machineCount
			case msg.MESSAGE:
				inMessage = message.(msg.MESSAGE).MSG
				messageChan <- inMessage
			}
		}
		leaderDown <- 1
		fmt.Println("--Sent leader down channel--")
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

func SendPaxosMessage(address string, message interface{}) {
	conn := Dial(address)
	if conn != nil {
		encoder := gob.NewEncoder(conn)
		//encoder := GetEncoder(address)
		//var err = encoder.Encode(&message)
		var err = encoder.Encode(&message)
		if err != nil {
			fmt.Println("WRONG!!! ", err)
		}
		/*if err != nil {
			fmt.Println("--TCP Paxos message: Encoding failed!!: ", err, "--")
			conn, err := net.Dial("tcp", address)
			if err != nil {
				fmt.Println(err)
			} else {
				encoder := gob.NewEncoder(conn)
				var enErr = encoder.Encode(&message)
				if enErr != nil {
					fmt.Println("Encoder problem: ", enErr)
				}
				StoreEncoder(conn, *encoder)
			}
		} else {
			Close(conn)
			StoreEncoder(conn, *encoder)
		}*/
	} else {
		fmt.Println("--Cannot send message to node--")
	}
}

func SendSomething(address string, message interface{}) {
	fmt.Println("1")
	conn := Dial(address)
	if conn != nil {
		fmt.Println("2")
		encoder := gob.NewEncoder(conn)
		//encoder := GetEncoder(address)
		//var err = encoder.Encode(&message)
		fmt.Println("3")
		var err = encoder.Encode(&message)
		fmt.Println("4")
		if err != nil {
			fmt.Println("WRONG!!! ", err)
		}
		/*if err != nil {
			fmt.Println("--TCP Paxos message: Encoding failed!!: ", err, "--")
			conn, err := net.Dial("tcp", address)
			if err != nil {
				fmt.Println(err)
			} else {
				encoder := gob.NewEncoder(conn)
				var enErr = encoder.Encode(&message)
				if enErr != nil {
					fmt.Println("Encoder problem: ", enErr)
				}
				StoreEncoder(conn, *encoder)
			}
		} else {
			Close(conn)
			StoreEncoder(conn, *encoder)
		}*/
	} else {
		fmt.Println("--Cannot send message to node--")
	}
	fmt.Println("5")
}

func AppendIfMissing(slice []node.T_Node, i node.T_Node) []node.T_Node {
	for _, ele := range slice {
		if ele.IP == i.IP {
			return slice
		}
	}
	return append(slice, i)
}
