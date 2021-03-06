package Paxos

import (
	"encoding/gob"
	"fmt"
	"lab5/Utils"
	"lab5/model/Network/message"
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
	fmt.Println("Paxos handler up!")
	service := "0.0.0.0:1338"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	Utils.CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	for {
		Utils.CheckError(err)
		conn, _ := listener.Accept()
		go holdConnection(conn)

	}
}

func holdConnection(conn net.Conn) {
	var connectionOK = true
	for connectionOK == true {
		decoder := gob.NewDecoder(conn)
		var msg interface{}
		err := decoder.Decode(&msg)
		if err != nil {
			connectionOK = false
		}
		if msg != nil {
			ip := Utils.GetIp(conn.RemoteAddr().String())
			switch msg.(type) {
			case message.Prepare:
				var mes = message.Wrapper{Ip: ip, Message: msg.(message.Prepare)}
				message.PrepareChan <- mes
			case message.Promise:
				var mes = message.Wrapper{Ip: ip, Message: msg.(message.Promise)}
				message.PromiseChan <- mes
			case message.Accept:
				var mes = message.Wrapper{Ip: ip, Message: msg.(message.Accept)}
				message.AcceptChan <- mes
			case message.Learn:
				var mes = message.Wrapper{Ip: ip, Message: msg.(message.Learn)}
				message.LearnChan <- mes
			}
		}
	}
	fmt.Println("Connection died, letting it go")
	fmt.Println(conn.RemoteAddr().String())
}
