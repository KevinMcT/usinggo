package px

import (
	"fmt"
	"lab5/model/net/tcp"
	"lab6/controller/node"
	"lab6/model/FifoList"
	"lab6/model/RoundVar"
	"lab6/model/net/msg"
	"time"
)

var (
	self           node.T_Node
	nodeList       = make([]node.T_Node, 0)
	round          int
	clientValue    interface{}
	promiseList    = make([]msg.Promise, 0)
	quorumPromise  bool
	waitPromisChan = make(chan string, 1)
	waiting        bool
	wmessages      = FifoList.NewQueue()
	newNodeOk      bool
)

func Proposer(me node.T_Node, nc chan node.T_Node, ac chan interface{}) {
	round = RoundVar.GetRound().Round
	self = me
	waiting = false
	quorumPromise = false
	newNodeOk = true
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

func handlePush(ac chan interface{}) {
	for {
		cv := <-ac
		wmessages.Add(cv)
	}
}

//Method for getting messages from the que and give them to
//paxos to be processed.
func handleMessages() {
	timeout := make(chan bool, 1)
	for {
		time.Sleep(25 * time.Millisecond)
		timeout <- true
		select {
		case nnAddress := <-msg.GetNodeOnTrack:
			fmt.Println("Getting node up to date!")
			newNodeOk = false
			sendAddress := nnAddress + ":1338"
			var prepare = msg.Prepare{ROUND: round}
			var updateMessage = msg.UpdateNode{PrepareMessage: prepare, SlotList: slots, BankAccounts: bankAccounts}
			var message interface{}
			message = updateMessage
			fmt.Println("Sending update Message")
			tcp.SendPaxosMessage(sendAddress, message)
		case <-timeout:
			if quorumPromise == true && newNodeOk == true {
				var msg = wmessages.Next()
				if msg != nil {
					clientValue = msg
					RoundVar.GetRound().MessageNumber = RoundVar.GetRound().MessageNumber + 1
					sendAccept()
				}
			}
		}
	}
}

func sendPrepare() {
	<-msg.SendPrepareChan
	quorumPromise = true
	actuallySendPrepare()
}

func actuallySendPrepare() {
	nodeList = RoundVar.GetRound().List
	RoundVar.GetRound().MessageNumber = lastLearntMsgNumber
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
		// If a node has come back up or a new one is added
		// we dont check for quorum.
		if newNodeOk == true {
			checkPromises()
		} else {
			fmt.Println("New node up to date, can continue!")
			newNodeOk = true
		}
	}
}

/*
When wait is finished we check to see if a value is in the promises.
*/
func checkPromises() {
	//allDefault := true
	if len(promiseList) > len(RoundVar.GetRound().List)/2 {
		quorumPromise = true
		promiseList = make([]msg.Promise, 0)
	}
}

/*
Method for choosing a value from all the promises.
*/
func pickValueFromProposeList() {
	var largestRound int = 0
	var largestRoundValue interface{} = nil
	for _, pMsg := range promiseList {
		if pMsg.LASTACCEPTEDROUND > largestRound {
			largestRound = pMsg.LASTACCEPTEDROUND
			largestRoundValue = pMsg.LASTACCEPTEDVALUE
		}
	}
	clientValue = largestRoundValue
	var stringMessage = fmt.Sprintf("Sending accept with value %s for round:%d messageNumber:%d ", clientValue, round, RoundVar.GetRound().MessageNumber)
	fmt.Println("--", stringMessage, "--")
	sendAccept()
	promiseList = make([]msg.Promise, 0)
}
