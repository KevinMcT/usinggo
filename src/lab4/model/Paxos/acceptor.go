package Paxos

import (
	"encoding/gob"
	"fmt"
	"lab4/model/Network/message"
	"lab4/model/RoundVar"
	"net"
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
		sendConn, err := net.Dial("tcp", sendAddress)
		if err == nil {
			encoder := gob.NewEncoder(sendConn)
			var learn = message.Learn{ROUND: promisedRound, VALUE: acceptedValue}
			var msg interface{}
			msg = learn
			encoder.Encode(&msg)
			sendConn.Close()
		} else {
			fmt.Println("Cannot send learn to node")
		}
	}
}

func sendPromise(address string) {
	address = address + ":1338"
	RoundVar.GetRound().Round = promisedRound
	fmt.Println("Promised to round: ", promisedRound)
	conn, err := net.Dial("tcp", address)
	if err == nil {
		encoder := gob.NewEncoder(conn)
		var promise = message.Promise{ROUND: promisedRound, LASTACCEPTEDROUND: lastAcceptedRound, LASTACCEPTEDVALUE: lastAcceptedValue}
		var msg interface{}
		msg = promise
		encoder.Encode(&msg)
		conn.Close()
	} else {
		fmt.Println("Cannot send promise to node")
	}
}
