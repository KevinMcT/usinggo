package px

import (
	//"encoding/gob"
	"fmt"
	"lab6/model/RoundVar"
	"lab6/model/net/msg"
	"lab6/model/net/tcp"
)

var (
	lastAcceptedValue     interface{}
	lastAcceptedRound     int
	lastAccpetedMsgNumber int

	promisedRound int
	acceptedValue interface{}
	hasAccepted   bool
	msgNumber     int
)

func Acceptor() {
	fmt.Println("--Proposer up and waiting ...")
	acceptedValue = "-1"
	lastAcceptedRound = -1
	lastAccpetedMsgNumber = -1
	lastAcceptedValue = "-1"
	promisedRound = 0
	msgNumber = 0
	go receivedPrepare()
	go receivedAccept()
}

func receivedPrepare() {
	for {
		value := <-msg.PrepareChan
		var tempRound = value.Message.(msg.Prepare).ROUND
		if tempRound >= promisedRound {
			promisedRound = value.Message.(msg.Prepare).ROUND
			sendPromise(value.Ip)
		}
	}
}

func receivedAccept() {
	for {
		value := <-msg.AcceptChan
		acceptMsg := value.Message.(msg.Accept)
		/*fmt.Println(promisedRound)
		fmt.Println(acceptMsg)*/
		if acceptMsg.ROUND == promisedRound {
			acceptedValue = acceptMsg.VALUE
			msgNumber = acceptMsg.MSGNUMBER
			lastAcceptedValue = acceptedValue
			lastAcceptedRound = promisedRound
			lastAccpetedMsgNumber = acceptMsg.MSGNUMBER
			sendLearn(value.Ip)
		} else {
			fmt.Println("--Error in number--")
		}
	}
}

func sendLearn(address string) {
	nodeList = RoundVar.GetRound().List
	for _, v := range nodeList {
		sendAddress := v.IP + ":1338"
		var learn = msg.Learn{ROUND: promisedRound, VALUE: acceptedValue, MSGNUMBER: msgNumber}
		var message interface{}
		message = learn
		tcp.SendPaxosMessage(sendAddress, message)
	}
}

func sendPromise(address string) {
	address = address + ":1338"
	RoundVar.GetRound().Round = promisedRound
	fmt.Println("--Promised to round: ", promisedRound, "--")
	var promise = msg.Promise{ROUND: promisedRound, LASTACCEPTEDROUND: lastAcceptedRound, LASTACCEPTEDVALUE: lastAcceptedValue}
	var message interface{}
	message = promise
	tcp.SendPaxosMessage(address, message)
}
