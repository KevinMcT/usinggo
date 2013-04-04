package px

import (
	//"encoding/gob"
	"fmt"
	"lab5Merge/controller/node"
	"lab5Merge/model/FifoList"
	"lab5Merge/model/RoundVar"
	"lab5Merge/model/net/msg"
	"lab5Merge/model/net/tcp"
	"sync"
	"time"
)

type WaitingMessages struct {
	mu       sync.Mutex
	messages *FifoList.Fifo
}

var (
	leader          node.T_Node
	self            node.T_Node
	nodeList        = make([]node.T_Node, 0)
	round           int
	msgNumber       int
	clientValue     string
	promiseList     = make([]msg.Promise, 0)
	quorumPromise   bool
	waitPromisChan  = make(chan string, 1)
	waiting         bool
	waitingMessages *WaitingMessages
)

func init() {
	waitingMessages = new(WaitingMessages)
	waitingMessages.messages = FifoList.NewFifo()
}

func Proposer(led node.T_Node, me node.T_Node, nc chan node.T_Node, ac chan string) {
	round = RoundVar.GetRound().Round
	leader = led
	self = me
	waiting = false
	quorumPromise = false
	go receviedPromise()
	go waitForPromise()
	if led.IP == me.IP {
		//fmt.Println(leader)
		//fmt.Println(self)
		go sendPrepare()
		go handleMessages()
	}
	for {
		cv := <-ac
		waitingMessages.messages.Push(cv)
	}
}

//Method for getting messages from the que and give them to
//paxos to be processed.
func handleMessages() {
	for {
		time.Sleep(2 * time.Second)
		var msg = waitingMessages.messages.Pop()
		clientValue = msg
		if quorumPromise == true {
			sendAccept()
			RoundVar.GetRound().MessageNumber = msgNumber + 1
		}
	}
}

func sendPrepare() {
	<-msg.SendPrepareChan
	fmt.Println("sending prepare!")
	quorumPromise = true
	nodeList = RoundVar.GetRound().List
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
		var accept = msg.Accept{ROUND: round, MSGNUMBER: msgNumber, VALUE: clientValue}
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
		time.Sleep(2000 * time.Millisecond)
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
	var stringMessage = fmt.Sprintf("Sending accept with value %s for round:%d messageNumber:%d ", clientValue, round, msgNumber)
	fmt.Println(stringMessage)
	sendAccept()
	promiseList = make([]msg.Promise, 0)
}

//Method for adding a incoming message from the client to the que of 
//waiting messages. 
func addMessageToQue(msg string) {
	waitingMessages.mu.Lock()
	waitingMessages.messages.Push(msg)
	waitingMessages.mu.Unlock()
}

//Message to get the last message from the que of waiting messages
func getNextMessageFromQue() string {
	waitingMessages.mu.Lock()
	var msg = waitingMessages.messages.Pop()
	waitingMessages.mu.Unlock()
	return msg
}
