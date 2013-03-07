package Paxos

import (
	"encoding/gob"
	"fmt"
	//lab4/Utils"
	"lab4/model/Network/message"
	"net"
)

var (
	leader      message.Node
	self        message.Node
	nodeList    = make([]message.Node, 0)
	round       int
	clientValue string
)

func Proposer(led message.Node, me message.Node, nc chan message.Node, ac chan string) {
	round = 0
	self = me
	go fillNodelist(nc)
	go receviedPromise()
	for {
		cv := <-ac
		clientValue = cv
		fmt.Println("Received value ", clientValue)
		round = round + 1
		sendPrepare()
	}
}

func fillNodelist(nc chan message.Node) {
	if self.LEAD == false {
		nodeList = append(nodeList, self)
	}
	if self.LEAD == true {
		leader = self
	}
	for {
		node := <-nc
		if node.LEAD == true {
			leader = node
		}
		nodeList = append(nodeList, node)
	}
}

func sendPrepare() {
	fmt.Println(nodeList)
	for _, v := range nodeList {
		//Send prepare
		sendAddress := v.IP + ":1338"
		fmt.Println("Sending prepare to ", sendAddress)
		sendConn, err := net.Dial("tcp", sendAddress)
		if err != nil {
			fmt.Println(err)
		} else {
			encoder := gob.NewEncoder(sendConn)
			var prepare = message.Prepare{ROUND: round}
			var msg interface{}
			msg = prepare
			encoder.Encode(&msg)
		}
		sendConn.Close()
	}
}

func receviedPromise() {
	for {
		value := <-message.PromiseChan
		promiseMsg := value.Message.(message.Promise)
		if promiseMsg.ROUND == round {
			address := value.Ip + ":1338"
			fmt.Println("Sending accept to ", address)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				fmt.Println(err)
			} else {
				encoder := gob.NewEncoder(conn)
				var accept = message.Accept{ROUND: round, VALUE: clientValue}
				var msg interface{}
				msg = accept
				encoder.Encode(&msg)
			}
			conn.Close()
		} else {
			fmt.Println("Node promised to higher round, not sending accept")
		}
	}
}
