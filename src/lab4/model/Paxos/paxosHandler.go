package Paxos

import (
	"encoding/gob"
	"fmt"
	"lab4/Utils"
	"lab4/model/Network/message"
	"net"
	"strings"
)

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
		conn, err := listener.Accept()

		decoder := gob.NewDecoder(conn)
		var msg interface{}
		err = decoder.Decode(&msg)
		if err != nil {
			fmt.Println(err)
		}
		conn.Close()
		if msg != nil {
			remote := conn.RemoteAddr().String()
			remoteSplit := strings.Split(remote, ":")
			switch msg.(type) {
			case message.Prepare:
				var mes = message.Wrapper{Ip: remoteSplit[0], Message: msg.(message.Prepare)}
				message.PrepareChan <- mes
			case message.Promise:
				var mes = message.Wrapper{Ip: remoteSplit[0], Message: msg.(message.Promise)}
				message.PromiseChan <- mes
			case message.Accept:
				var mes = message.Wrapper{Ip: remoteSplit[0], Message: msg.(message.Accept)}
				message.AcceptChan <- mes
			case message.Learn:
				var mes = message.Wrapper{Ip: remoteSplit[0], Message: msg.(message.Learn)}
				message.LearnChan <- mes
			}
		}
	}
}
