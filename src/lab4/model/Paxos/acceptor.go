package Paxos

import (
	"encoding/gob"
	"fmt"
	"lab4/model/Network/message"
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
		fmt.Println("Received Prepare for round ", value.Message.(message.Prepare).ROUND)
		promisedRound = value.Message.(message.Prepare).ROUND
		sendPromise(value.Ip)
	}
}

func receivedAccept() {
	for {
		value := <-message.AcceptChan
		fmt.Println("Received accept!")
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
	for _, v := range nodeList {
		sendAddress := v.IP + ":1338"
		fmt.Println("Sending learn to ", sendAddress)
		sendConn, err := net.Dial("tcp", sendAddress)
		if err != nil {
			fmt.Println(err)
		} else {
			encoder := gob.NewEncoder(sendConn)
			var learn = message.Learn{ROUND: round, VALUE: acceptedValue}
			var msg interface{}
			msg = learn
			encoder.Encode(&msg)
		}
		sendConn.Close()
	}
}

func sendPromise(address string) {
	address = address + ":1338"
	fmt.Println("Sending promise to ", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
	} else {
		encoder := gob.NewEncoder(conn)
		var promise = message.Promise{ROUND: promisedRound, LASTACCEPTEDROUND: lastAcceptedRound, LASTACCEPTEDVALUE: lastAcceptedValue}
		var msg interface{}
		msg = promise
		encoder.Encode(&msg)
	}
	conn.Close()
}
