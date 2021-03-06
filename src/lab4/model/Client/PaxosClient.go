package main

import (
	"encoding/gob"
	"fmt"
	"lab4/model/Network/message"
	"net"
)

/*
The paxos client used to communicate with the paxos system.
User is asked for a ip and a one word message to send over 
to the system. If the entered ip is not correct/listning the user
is prompted to enter a new one
*/
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
		if err == nil {
			fmt.Println("Enter a value to send")
			var st string
			fmt.Scanf("%s", &st)
			defer conn.Close()
			encoder := gob.NewEncoder(conn)

			var sendMsg = message.ClientRequestMessage{Content: st}
			var msg interface{}
			msg = sendMsg

			encoder.Encode(&msg)
			fmt.Println("Message sent to paxos replica")
		} else {
			fmt.Println("Seems like the node you are trying to connect is gone down or does not exist. Please try another address")
		}
	}
}
