package px

import (
	//"encoding/gob"
	"fmt"
	"lab5Merge/model/RoundVar"
	"lab5Merge/model/net/msg"
	"lab5Merge/model/net/tcp"
)

var (
	lastAcceptedValue     string
	lastAcceptedRound     int
	lastAccpetedMsgNumber int

	promisedRound int
	acceptedValue string
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
		promisedRound = value.Message.(msg.Prepare).ROUND
		sendPromise(value.Ip)
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
