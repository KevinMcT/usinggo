package Paxos

import (
	"encoding/gob"
	"fmt"
	"lab5/model/FifoList"
	"lab5/model/Network/message"
	"lab5/model/Network/tcp"
	"lab5/model/RoundVar"
	"sync"
	"time"
)

type WaitingMessages struct {
	mu       sync.Mutex
	messages *FifoList.Fifo
}

var (
	leader          message.Node
	self            message.Node
	nodeList        = make([]message.Node, 0)
	round           int
	msgNumber       int
	clientValue     string
	promiseList     = make([]message.Promise, 0)
	quorumPromise   bool
	waitPromisChan  = make(chan string, 1)
	waiting         bool
	waitingMessages *WaitingMessages
)

func init() {
	waitingMessages = new(WaitingMessages)
	waitingMessages.messages = FifoList.NewFifo()
}

func Proposer(led message.Node, me message.Node, nc chan message.Node, ac chan string) {
	round = RoundVar.GetRound().Round
	self = me
	waiting = false
	quorumPromise = false
	go receviedPromise()
	go waitForPromise()
	if led.IP == me.IP {
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
	<-message.SendPrepareChan
	quorumPromise = true
	nodeList = RoundVar.GetRound().List
	for _, v := range nodeList {
		sendAddress := v.IP + ":1338"
		sendConn := tcp.Dial(sendAddress)
		if sendConn != nil {
			encoder := gob.NewEncoder(sendConn)
			var prepare = message.Prepare{ROUND: round}
			var msg interface{}
			msg = prepare
			encoder.Encode(&msg)
			tcp.Close(sendConn)
		} else {
			fmt.Println("Cannot send prepare to node")
		}
	}
}

func sendAccept() {
	nodeList = RoundVar.GetRound().List
	for _, v := range nodeList {
		address := v.IP + ":1338"
		conn := tcp.Dial(address)
		if conn != nil {
			encoder := gob.NewEncoder(conn)
			msgNumber = RoundVar.GetRound().MessageNumber
			var accept = message.Accept{ROUND: round, MSGNUMBER: msgNumber, VALUE: clientValue}
			var msg interface{}
			msg = accept
			var err = encoder.Encode(&msg)
			if err != nil {
				fmt.Println("Encoding failed!!: ", err)
			}
			tcp.Close(conn)
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

/*
After the first promise has been received
we start a timer to wait for the rest.
*/
func waitForPromise() {
	for {
		<-waitPromisChan
		waiting = true
		time.Sleep(3 * time.Second)
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
	promiseList = make([]message.Promise, 0)
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
