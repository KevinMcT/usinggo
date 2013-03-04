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
)

func Proposer() {
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
		fmt.Println("Received Prepare!!")
		sendPromise(value.Ip)
	}
}

func receivedAccept() {
	for {
		value := <-message.AcceptChan
		fmt.Println("Received accept!")
		sendLearn(value.Ip)
	}
}

func sendLearn(address string) {
	address = address + ":1338"
	fmt.Println("Sending learn to ", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
	} else {
		encoder := gob.NewEncoder(conn)
		var promise = message.Learn{ROUND: 1, VALUE: "test"}
		var msg interface{}
		msg = promise
		encoder.Encode(&msg)
	}
	conn.Close()
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
