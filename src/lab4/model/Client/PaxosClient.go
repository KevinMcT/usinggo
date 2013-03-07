package main

import (
	"encoding/gob"
	"fmt"
	"lab4/Utils"
	"lab4/model/Network/message"
	"net"
)

func main() {
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
		Utils.CheckError(err)
		defer conn.Close()
		encoder := gob.NewEncoder(conn)

		var sendMsg = message.ClientRequestMessage{Content: "Hello world!"}
		var msg interface{}
		msg = sendMsg

		encoder.Encode(&msg)
		fmt.Println("Message sent to paxos replica")
	}
}
