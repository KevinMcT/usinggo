package px

import (
	//"encoding/gob"
	"fmt"
	"lab5Merge/controller/node"
	"lab5Merge/model/FifoList"
	"lab5Merge/model/RoundVar"
	"lab5Merge/model/net/msg"
	"lab5Merge/model/net/tcp"
	//"sync"
	"time"
)

var (
	self           node.T_Node
	nodeList       = make([]node.T_Node, 0)
	round          int
	clientValue    string
	promiseList    = make([]msg.Promise, 0)
	quorumPromise  bool
	waitPromisChan = make(chan string, 1)
	waiting        bool
	wmessages      = FifoList.NewQueue()
)

func Proposer(me node.T_Node, nc chan node.T_Node, ac chan string) {
	round = RoundVar.GetRound().Round
	self = me
	waiting = false
	quorumPromise = false
	go receviedPromise()
	go waitForPromise()
	go handlePush(ac)
	go handleMessages()
	if RoundVar.GetRound().CurrentLeader.IP == me.IP {
		sendPrepare()
	}
	for {
		<-msg.RestartProposer
		if RoundVar.GetRound().CurrentLeader.IP == me.IP {
			round = RoundVar.GetRound().Round
			actuallySendPrepare()
		}
	}

}

func handlePush(ac chan string) {
	for {
		cv := <-ac
		wmessages.Add(cv)
	}
}

//Method for getting messages from the que and give them to
//paxos to be processed.
func handleMessages() {
	for {
		time.Sleep(25 * time.Millisecond)
		//var msg = waitingMessages.messages.Pop()
		var msg = wmessages.Next()
		//fmt.Println("--Pop'ed message from queue")
		if msg != nil {
			clientValue = msg.(string)
			if quorumPromise == true {
				RoundVar.GetRound().MessageNumber = RoundVar.GetRound().MessageNumber + 1
				sendAccept()
			}
		}
	}
}

func sendPrepare() {
	<-msg.SendPrepareChan
	fmt.Println("sending prepare!")
	quorumPromise = true
	actuallySendPrepare()
}

func actuallySendPrepare() {
	nodeList = RoundVar.GetRound().List
	RoundVar.GetRound().MessageNumber = lastLearntMsgNumber
	fmt.Println(RoundVar.GetRound().MessageNumber)
	for _, v := range nodeList {
		sendAddress := v.IP + ":1338"
		var prepare = msg.Prepare{ROUND: round}
		var message interface{}
		message = prepare
		tcp.SendPaxosMessage(sendAddress, message)
	}
}

func sendAccept() {
	nodeList = RoundVar.GetRound().List
	for _, v := range nodeList {
		address := v.IP + ":1338"
		var accept = msg.Accept{ROUND: round, MSGNUMBER: RoundVar.GetRound().MessageNumber, VALUE: clientValue}
		var message interface{}
		message = accept
		tcp.SendPaxosMessage(address, message)
	}
}

func receviedPromise() {
	for {
		value := <-msg.PromiseChan
		if waiting == false {
			waitPromisChan <- "wait"
		}
		promiseList = append(promiseList, value.Message.(msg.Promise))
	}
}

/*
After the first promise has been received
we start a timer to wait for the rest.
*/
func waitForPromise() {
	for {
		<-waitPromisChan
		waiting = true
		time.Sleep(10 * time.Millisecond)
		waiting = false
		checkPromises()
	}
}

/*
When wait is finished we check to see if a value is in the promises.
*/
func checkPromises() {
	//allDefault := true
	if len(promiseList) > len(RoundVar.GetRound().List)/2 {
		quorumPromise = true
		/*for _, pMsg := range promiseList {
			if pMsg.LASTACCEPTEDVALUE != "-1" {
				allDefault = false
			}
		}
		if allDefault == true {
			sendAccept()
			promiseList = make([]message.Promise, 0)
		} else {
			pickValueFromProposeList()
		}*/
	}
}

/*
Method for choosing a value from all the promises.
*/
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
	var stringMessage = fmt.Sprintf("Sending accept with value %s for round:%d messageNumber:%d ", clientValue, round, RoundVar.GetRound().MessageNumber)
	fmt.Println(stringMessage)
	sendAccept()
	promiseList = make([]msg.Promise, 0)
}
