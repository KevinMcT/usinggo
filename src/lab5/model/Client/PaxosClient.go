package main

import (
	"encoding/gob"
	"fmt"
	"lab5/Utils"
	"lab5/model/Network/message"
	"net"
	"time"
)

/*
The paxos client used to communicate with the paxos system.
User is asked for a ip and a one word message to send over 
to the system. If the entered ip is not correct/listning the user
is prompted to enter a new one
*/
func main() {
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
			//for i := 0; i < 10; i++ {
			encoder := gob.NewEncoder(conn)
			//var stringMessage = fmt.Sprintf("%s%d", st, i)			
			var sendMsg = message.ClientRequestMessage{Content: st}
			var msg interface{}
			msg = sendMsg
			encoder.Encode(&msg)
			fmt.Println("Message sent to paxos replica")
			time.Sleep(500 * time.Millisecond)
			//}
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
		go holdConnection(conn)
	}
}

func holdConnection(conn net.Conn) {
	var connectionOK = true
	for connectionOK == true {
		fmt.Println("Got response in client")
		decoder := gob.NewDecoder(conn)
		var msg interface{}
		err := decoder.Decode(&msg)
		if err != nil {
			connectionOK = false
		}
		if msg != nil {
			var clientMsg message.ClientResponseMessage
			clientMsg = msg.(message.ClientResponseMessage)
			fmt.Println("---------------------------------------------------")
			fmt.Println(clientMsg.Content)
			fmt.Println("---------------------------------------------------")
		} else {
			fmt.Println("Message is empty stupid!")
		}
	}
	fmt.Println("Client closed connection, no more to share")
}
