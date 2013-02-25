package Client

import (
	"encoding/gob"
	"fmt"
	"lab4/model/Network/message"
	"net"
	"os"
)

func PaxosClient() {
	ConnectToPaxos()

}

func ConnectToPaxos() {
	fmt.Println("Enter ip to connecto to")
	var ip string
	fmt.Scanf("%s", &ip)
	fmt.Println("Connecting to Paxos replica")
	service := ip + ":1337"
	fmt.Println(service)
	conn, err := net.Dial("tcp", service)
	checkError(err)
	defer conn.Close()
	encoder := gob.NewEncoder(conn)
	var msg = message.ClientRequestMessage{"content"}
	encoder.Encode(msg)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error", err.Error())
		os.Exit(1)
	}
}
