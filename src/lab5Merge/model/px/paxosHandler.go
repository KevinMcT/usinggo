package px

import (
	"encoding/gob"
	"fmt"
	"lab5Merge/Utils"
	"lab5Merge/model/net/msg"
	//"lab5Merge/model/net/tcp"
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
	decoder := gob.NewDecoder(conn)
	for connectionOK == true {
		var message interface{}
		err := decoder.Decode(&message)
		if err != nil {
			connectionOK = false
		}
		if message != nil {
			ip := Utils.GetIp(conn.RemoteAddr().String())
			switch message.(type) {
			case msg.Prepare:
				fmt.Println("received Prepare")
				var mes = msg.Wrapper{Ip: ip, Message: message.(msg.Prepare)}
				msg.PrepareChan <- mes
			case msg.Promise:
				fmt.Println("received Promise")
				var mes = msg.Wrapper{Ip: ip, Message: message.(msg.Promise)}
				msg.PromiseChan <- mes
			case msg.Accept:
				fmt.Println("received Accept")
				var mes = msg.Wrapper{Ip: ip, Message: message.(msg.Accept)}
				msg.AcceptChan <- mes
			case msg.Learn:
				fmt.Println("received Learn")
				var mes = msg.Wrapper{Ip: ip, Message: message.(msg.Learn)}
				msg.LearnChan <- mes
			}
		}
	}
	//tcp.StoreDecoder(conn, *decoder)
	fmt.Println("Connection died, letting it go")
	fmt.Println(conn.RemoteAddr().String())
}
