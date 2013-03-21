package Paxos

import (
	"encoding/gob"
	"fmt"
	"lab5/model/Network/message"
	"lab5/model/Network/tcp"
	"lab5/model/RoundVar"
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
	go receivedPrepare()
	go receivedAccept()
}

func receivedPrepare() {
	for {
		value := <-message.PrepareChan
		promisedRound = value.Message.(message.Prepare).ROUND
		sendPromise(value.Ip)
	}
}

func receivedAccept() {
	for {
		value := <-message.AcceptChan
		//fmt.Println("received accept")
		acceptMsg := value.Message.(message.Accept)
		if acceptMsg.ROUND == promisedRound {
			acceptedValue = acceptMsg.VALUE
			lastAcceptedValue = acceptedValue
			lastAcceptedRound = promisedRound
			sendLearn(value.Ip)
		}
	}
}

func sendLearn(address string) {
	nodeList = RoundVar.GetRound().List
	for _, v := range nodeList {
		sendAddress := v.IP + ":1338"
		sendConn := tcp.Dial(sendAddress)
		if sendConn != nil {
			encoder := gob.NewEncoder(sendConn)
			var learn = message.Learn{ROUND: promisedRound, VALUE: acceptedValue}
			var msg interface{}
			msg = learn
			encoder.Encode(&msg)
			//fmt.Println("Sending learn")
			tcp.Close(sendConn)
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
		encoder := gob.NewEncoder(conn)
		var promise = message.Promise{ROUND: promisedRound, LASTACCEPTEDROUND: lastAcceptedRound, LASTACCEPTEDVALUE: lastAcceptedValue}
		var msg interface{}
		msg = promise
		encoder.Encode(&msg)
		//fmt.Println("Sending promise")
		tcp.Close(conn)
	} else {
		fmt.Println("Cannot send promise to node")
	}
}
