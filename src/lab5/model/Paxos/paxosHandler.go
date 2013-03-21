package Paxos

import (
	"encoding/gob"
	"fmt"
	"lab5/Utils"
	"lab5/model/Network/message"
	"net"
	"strings"
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
		//fmt.Println("Waiting again....")
		conn, _ := listener.Accept()
		go holdConnection(conn)

	}
}

func holdConnection(conn net.Conn) {
	var connectionOK = true
	for connectionOK == true {
		//fmt.Println("Received something!")
		decoder := gob.NewDecoder(conn)
		var msg interface{}
		err := decoder.Decode(&msg)
		if err != nil {
			connectionOK = false
		}
		if msg != nil {
			remote := conn.RemoteAddr().String()
			remoteSplit := strings.Split(remote, ":")
			switch msg.(type) {
			case message.Prepare:
				//fmt.Println("Received prepare")
				var mes = message.Wrapper{Ip: remoteSplit[0], Message: msg.(message.Prepare)}
				message.PrepareChan <- mes
			case message.Promise:
				//fmt.Println("Received promise")
				var mes = message.Wrapper{Ip: remoteSplit[0], Message: msg.(message.Promise)}
				message.PromiseChan <- mes
			case message.Accept:
				//fmt.Println("Received accept")
				var mes = message.Wrapper{Ip: remoteSplit[0], Message: msg.(message.Accept)}
				message.AcceptChan <- mes
			case message.Learn:
				//fmt.Println("Received learn")
				var mes = message.Wrapper{Ip: remoteSplit[0], Message: msg.(message.Learn)}
				message.LearnChan <- mes
			}
		}
	}
	fmt.Println("Connection died, letting it go")
}
