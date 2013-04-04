package px

import (
	//"encoding/gob"
	"fmt"
	"lab5Merge/model/RoundVar"
	"lab5Merge/model/net/msg"
	"lab5Merge/model/net/tcp"
)

var (
	lastAcceptedValue string
	lastAcceptedRound int
	promisedRound     int
	acceptedValue     string
	hasAccepted       bool
)

func Acceptor() {
	fmt.Println("Proposer up and waiting ...")
	acceptedValue = "-1"
	lastAcceptedRound = -1
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
		if acceptMsg.ROUND == promisedRound {
			acceptedValue = acceptMsg.VALUE
			msgNumber = acceptMsg.MSGNUMBER
			lastAcceptedValue = acceptedValue
			lastAcceptedRound = promisedRound
			sendLearn(value.Ip)
		} else {
			fmt.Println("Wrong round number bitch!")
			fmt.Println(acceptMsg.ROUND)
		}
	}
}

func sendLearn(address string) {
	nodeList = RoundVar.GetRound().List
	for _, v := range nodeList {
		sendAddress := v.IP + ":1338"
		sendConn := tcp.Dial(sendAddress)
		if sendConn != nil {
			//encoder := gob.NewEncoder(sendConn)
			encoder := tcp.GetEncoder(sendAddress)
			var learn = msg.Learn{ROUND: promisedRound, VALUE: acceptedValue, MSGNUMBER: msgNumber}
			var message interface{}
			message = learn
			encoder.Encode(&message)
			//fmt.Println("Sending learn")
			tcp.Close(sendConn)
			tcp.StoreEncoder(sendConn, *encoder)
		} else {
			fmt.Println("Cannot send learn to node")
		}
	}
}

func sendPromise(address string) {
	address = address + ":1338"
	RoundVar.GetRound().Round = promisedRound
	fmt.Println("Promised to round: ", promisedRound)
	conn := tcp.Dial(address)
	if conn != nil {
		//encoder := gob.NewEncoder(conn)
		encoder := tcp.GetEncoder(address)
		var promise = msg.Promise{ROUND: promisedRound, LASTACCEPTEDROUND: lastAcceptedRound, LASTACCEPTEDVALUE: lastAcceptedValue}
		var message interface{}
		message = promise
		encoder.Encode(&message)
		//fmt.Println("Sending promise")
		tcp.Close(conn)
		tcp.StoreEncoder(conn, *encoder)
	} else {
		fmt.Println("Cannot send promise to node")
	}
}
