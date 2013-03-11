package Paxos

import (
	"encoding/gob"
	"fmt"
	//lab4/Utils"
	"lab4/model/Network/message"
	"lab4/model/RoundVar"
	"net"
	"time"
)

type Pair struct {
	Key, Value string
}

var (
	leader         message.Node
	self           message.Node
	nodeList       = make([]message.Node, 0)
	round          int
	clientValue    string
	promiseList    = make([]message.Promise, 0)
	waitPromisChan = make(chan string, 1)
	waiting        bool
)

func Proposer(led message.Node, me message.Node, nc chan message.Node, ac chan string) {
	round = RoundVar.GetRound().Round
	self = me
	waiting = false
	go receviedPromise()
	go waitForPromise()
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
	}
}

func sendPrepare() {
	nodeList = RoundVar.GetRound().List
	for _, v := range nodeList {
		//Send prepare
		sendAddress := v.IP + ":1338"
		sendConn, err := net.Dial("tcp", sendAddress)
		if err == nil {
			encoder := gob.NewEncoder(sendConn)
			var prepare = message.Prepare{ROUND: round}
			var msg interface{}
			msg = prepare
			encoder.Encode(&msg)
			sendConn.Close()
		} else {
			fmt.Println("Cannot send prepare to node")
		}
	}
}

func sendAccept() {
	nodeList = RoundVar.GetRound().List
	for _, v := range nodeList {
		address := v.IP + ":1338"
		conn, err := net.Dial("tcp", address)
		if err == nil {
			encoder := gob.NewEncoder(conn)
			var accept = message.Accept{ROUND: round, VALUE: clientValue}
			var msg interface{}
			msg = accept
			encoder.Encode(&msg)
			conn.Close()
		} else {
			fmt.Println("Cannot send accept to node")
		}
	}
}

func receviedPromise() {
	for {
		value := <-message.PromiseChan
		if waiting == false {
			waitPromisChan <- "wait"
		}
		promiseList = append(promiseList, value.Message.(message.Promise))
	}
}

func waitForPromise() {
	for {
		<-waitPromisChan
		waiting = true
		time.Sleep(2 * time.Second)
		waiting = false
		fmt.Println(promiseList)
		checkPromises()
	}
}

func checkPromises() {
	allDefault := true
	allNotDefault := true
	for _, pMsg := range promiseList {
		if pMsg.LASTACCEPTEDVALUE != "-1" {
			allDefault = false
		}
		if pMsg.LASTACCEPTEDVALUE == "-1" {
			allNotDefault = false
		}
	}
	if allDefault == true || allNotDefault == true {
		sendAccept()
		promiseList = make([]message.Promise, 0)
	} else {
		pickValueFromProposeList()
	}
}

func pickValueFromProposeList() {
	var largestRound int = 0
	var largestRoundValue string = ""
	for _, pMsg := range promiseList {
		if pMsg.LASTACCEPTEDROUND > largestRound {
			largestRound = pMsg.LASTACCEPTEDROUND
			largestRoundValue = pMsg.LASTACCEPTEDVALUE
		}
	}
	clientValue = largestRoundValue
	sendAccept()
	promiseList = make([]message.Promise, 0)
}
