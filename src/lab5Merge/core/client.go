package core

import (
	"encoding/gob"
	"fmt"
	"lab5Merge/Utils"
	"lab5Merge/controller/node"
	"lab5Merge/model/net/msg"
	"net"
	"os"
	"time"
)

var (
	serverList []node.T_Node
	sendAll    bool
	sentChan   = make(chan int, 1)
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
		fmt.Println("Waiting for servers... This might take up to 5 seconds, but you can still send a message")
		go GetServers(conn, GetIP())
		if err == nil {
			fmt.Println("Enter a value to send")
			var st string
			fmt.Scanf("%s", &st)
			var all string
			fmt.Println("Wait for response before sending new message? Y/N")
			fmt.Scanf("%s", &all)
			if all == "Y" || all == "y" || all == "yes" {
				sendAll = false
				sendToPaxos(st, conn)
			}
			if all == "N" || all == "n" || all == "no" {
				sendAll = true
				sendToPaxos(st, conn)
			}

		} else {
			fmt.Println("Seems like the node you are trying to connect is gone down or does not exist. Please try another address")
		}
	}
}

func sendToPaxos(st string, conn net.Conn) {
	for i := 0; i < 200; i++ {
		encoder := gob.NewEncoder(conn)
		var stringMessage = fmt.Sprintf("%s%d", st, i)
		var sendMsg = msg.ClientRequestMessage{Content: stringMessage}
		var message interface{}
		message = sendMsg
		encoder.Encode(&message)
		if sendAll == false {
			<-sentChan
		}
		time.Sleep(10 * time.Millisecond)
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
		conn, _ := listener.Accept()
		go holdClientConnection(conn)
	}
}

func holdClientConnection(conn net.Conn) {
	var connectionOK = true
	decoder := gob.NewDecoder(conn)
	for connectionOK == true {
		var message interface{}
		err := decoder.Decode(&message)
		switch message.(type) {
		case msg.ClientResponseMessage:
			if err != nil {
				connectionOK = false
			}
			if message != nil {
				var clientMsg msg.ClientResponseMessage
				clientMsg = message.(msg.ClientResponseMessage)
				fmt.Println("---------------------------------------------------")
				fmt.Println(clientMsg.Content)
				fmt.Println("---------------------------------------------------")
				if sendAll == false {
					sentChan <- 1
				}
			} else {
				fmt.Println("Message is empty stupid!")
			}
		case msg.ClientResponseNodes:
			var clientMsg msg.ClientResponseNodes
			clientMsg = message.(msg.ClientResponseNodes)
			serverList = clientMsg.List
		}
	}
	fmt.Println("Paxos closed connection, no more to share")
}

func GetServers(conn net.Conn, myIP string) {
	for {
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(5000 * time.Millisecond)
			timeout <- true
		}()
		select {
		case <-timeout:
			encoder := gob.NewEncoder(conn)
			var sendMsg = msg.ClientRequestNodes{RemoteAddress: myIP + ":1337"}
			var message interface{}
			message = sendMsg
			encoder.Encode(&message)
		}
	}
}

func GetIP() string {
	name, _ := os.Hostname()
	addr, _ := net.LookupHost(name)
	UDPAddr, _ := net.ResolveUDPAddr("udp4", addr[0]+":1888")
	ip := UDPAddr.IP.String()
	return ip
}
