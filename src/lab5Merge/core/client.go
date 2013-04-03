package core

import (
	"encoding/gob"
	"fmt"
	"lab5Merge/Utils"
	"lab5Merge/model/net/msg"
	"net"
	"time"
)

/*
The paxos client used to communicate with the paxos system.
User is asked for a ip and a one word message to send over 
to the system. If the entered ip is not correct/listning the user
is prompted to enter a new one
*/
func Client() {
	go waitForResponse()
	ConnectToPaxos()
}

func ConnectToPaxos() {

	for {
		fmt.Println("Enter ip to connecto to")
		var ip string
		fmt.Scanf("%s", &ip)

		fmt.Println("Connecting to Paxos replica")
		service := ip + ":1337"
		fmt.Println(service)
		conn, err := net.Dial("tcp", service)
		if err == nil {
			fmt.Println("Enter a value to send")
			var st string
			fmt.Scanf("%s", &st)
			for i := 0; i < 10; i++ {
				encoder := gob.NewEncoder(conn)
				var stringMessage = fmt.Sprintf("%s%d", st, i)
				var sendMsg = msg.ClientRequestMessage{Content: stringMessage}
				var message interface{}
				message = sendMsg
				encoder.Encode(&message)
				//fmt.Println("Message sent to paxos replica")
				time.Sleep(1000 * time.Millisecond)
			}
		} else {
			fmt.Println("Seems like the node you are trying to connect is gone down or does not exist. Please try another address")
		}
		conn.Close()
	}
}

func waitForResponse() {
	fmt.Println("Waiting for request responses")
	service := "0.0.0.0:1337"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	Utils.CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	for {
		Utils.CheckError(err)
		fmt.Println("Waiting for response in client")
		conn, _ := listener.Accept()
		go holdClientConnection(conn)
	}
}

func holdClientConnection(conn net.Conn) {
	var connectionOK = true
	for connectionOK == true {
		decoder := gob.NewDecoder(conn)
		var message interface{}
		err := decoder.Decode(&message)
		if err != nil {
			connectionOK = false
		}
		if message != nil {
			var clientMsg msg.ClientResponseMessage
			clientMsg = message.(msg.ClientResponseMessage)
			fmt.Println("---------------------------------------------------")
			fmt.Println(clientMsg.Content)
			fmt.Println("---------------------------------------------------")
		} else {
			fmt.Println("Message is empty stupid!")
		}
	}
	fmt.Println("Client closed connection, no more to share")
}
