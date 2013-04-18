package px

import (
	"encoding/gob"
	"fmt"
	"lab6/Utils"
	"lab6/model/net/msg"
	"net"
)

/*
Class to handle paxosmessages coming inn over tcp. When
a message is received it is checked and sendt out
on the appropriat channel. 
*/

var ()

func PaxosHandler() {
	HandlePaxosMessages()
}

func HandlePaxosMessages() {
	fmt.Println("--Paxos handler up!--")
	service := "0.0.0.0:1338"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	Utils.CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	for {
		Utils.CheckError(err)
		fmt.Println("PH waiting for connection!")
		conn, _ := listener.Accept()
		fmt.Println("PH got connection!")
		go holdConnection(conn)
	}
}

func holdConnection(conn net.Conn) {
	var connectionOK = true
	decoder := gob.NewDecoder(conn)
	for connectionOK == true {
		var message interface{}
		err := decoder.Decode(&message)
		if err != nil {
			connectionOK = false
		}
		if message != nil {
			fmt.Println("Got a new mesage!")
			ip := Utils.GetIp(conn.RemoteAddr().String())
			switch message.(type) {
			case msg.Prepare:
				var mes = msg.Wrapper{Ip: ip, Message: message.(msg.Prepare)}
				msg.PrepareChan <- mes
			case msg.Promise:
				var mes = msg.Wrapper{Ip: ip, Message: message.(msg.Promise)}
				msg.PromiseChan <- mes
			case msg.Accept:
				var mes = msg.Wrapper{Ip: ip, Message: message.(msg.Accept)}
				msg.AcceptChan <- mes
			case msg.Learn:
				var mes = msg.Wrapper{Ip: ip, Message: message.(msg.Learn)}
				msg.LearnChan <- mes
			case msg.UpdateNode:
				fmt.Println("Received update message!")
				var update = message.(msg.UpdateNode)
				slots = update.SlotList
				fmt.Println("SlotList: ", slots)
				bankAccounts = update.BankAccounts
				fmt.Println("BankAccounts: ", bankAccounts)
				var prepare = update.PrepareMessage
				var mes = msg.Wrapper{Ip: ip, Message: prepare}
				fmt.Println("Prepare: ", prepare)
				msg.PrepareChan <- mes
			}
		}
	}
	//tcp.StoreDecoder(conn, *decoder)
}
