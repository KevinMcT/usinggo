package core

import (
	"encoding/gob"
	"fmt"
	"lab6/Utils"
	"lab6/controller/node"
	"lab6/model/net/msg"
	"lab6/model/net/udp"
	"net"
	"os"
	"time"
)

var (
	serverListBank   []node.T_Node
	sendAllBank      bool
	sentChanBank     = make(chan int, 1)
	paxosAddressBank string

	lastConfirmedValueBank     string
	lastConfirmedMsgNumberBank int
	udpListenChanBank          = make(chan string, 0)
	bankMessage                interface{}
	connectionIP               string
	testChan                   = make(chan int, 0)
	conn1                      net.Conn
)

/*
The paxos client used to communicate with the paxos system.
User is asked for a ip and a one word message to send over 
to the system. If the entered ip is not correct/listning the user
is prompted to enter a new one
*/
func Bank() {
	go udp.Listen(udpListenChanBank)
	go waitForResponseBank()
	ConnectToPaxosBank()
}

func ConnectToPaxosBank() {
	udp.SendLocator(GetIPBank())
	ip := <-udpListenChanBank
	connectionIP = ip
	sendAllBank = true
	fmt.Println("Connecting to Paxos replica")
	service := connectionIP + ":1337"
	paxosAddressBank = service
	conn, err := net.Dial("tcp", paxosAddressBank)
	conn1 = conn
	go GetServersBank(conn1, GetIPBank())
	fmt.Println("Waiting for servers... This might take up to 5 seconds, but you can still send a message")
	for {
		lastConfirmedValueBank = ""
		lastConfirmedMsgNumberBank = -1
		if err == nil {
			fmt.Println("*********************************************************************")
			fmt.Println("*          Welcome to the National Bank of Bullshit                 *")
			fmt.Println("* Chconn1oose the function you would like to test using numbers 1 - 2    *")
			fmt.Println("* 1 - Deposit                                                       *")
			fmt.Println("* 2 - Withdraw                                                      *")
			fmt.Println("* 3 - Transfer                                                      *")
			fmt.Println("* 4 - Balance                                                       *")
			fmt.Println("* 0 - Quit                                                          *")
			fmt.Println("*********************************************************************")
			var in int
			fmt.Scanf("%d", &in)
			var accFrm string
			var accTo string
			var amt int
			switch in {
			case 1:
				fmt.Println("-- Enter your account number --")
				fmt.Scanf("%s", &accFrm)
				fmt.Println("-- Enter amount you wish to deposit --")
				fmt.Scanf("%d", &amt)
				bankMessage = msg.Deposit{AccountNumber: accFrm, Amount: amt}
				break
			case 2:
				fmt.Println("-- Enter your account number --")
				fmt.Scanf("%s", &accFrm)
				fmt.Println("-- Enter amount you wish to withdraw --")
				fmt.Scanf("%d", &amt)
				bankMessage = msg.Withdraw{AccountNumber: accFrm, Amount: amt}
				break
			case 3:
				fmt.Println("-- Enter your account number --")
				fmt.Scanf("%s", &accFrm)
				fmt.Println("-- Enter amount you wish to transfer --")
				fmt.Scanf("%d", &amt)
				fmt.Println("-- Enter account to transfer to --")
				fmt.Scanf("%s", &accTo)
				bankMessage = msg.Transfer{FromAccount: accFrm, ToAccount: accTo, Amount: amt}
				break
			case 4:
				fmt.Println("-- Enter your account number --")
				fmt.Scanf("%s", &accFrm)
				bankMessage = msg.Balance{AccountNumber: accFrm}
				break
			case 0:
				os.Exit(0)
			}
			sendToPaxosBank(bankMessage, conn1)
		} else {
			//fmt.Println("--Seems like the node you are trying to connect is gone down or does not exist. Please try another address--")
			//for i := 0; i < len(serverListBank); i++ {
			//	if serverListBank[i].IP != connectionIP {
			//		connectionIP = serverListBank[i].IP
			//	}
			//}
		}
	}
}

func sendToPaxosBank(st interface{}, conn net.Conn) {
	encoder := gob.NewEncoder(conn1)
	var sendMsg = msg.ClientRequestMessage{Content: st}
	var message interface{}
	message = sendMsg
	var err = encoder.Encode(&message)
	if err != nil {
		fmt.Println("Finding new address")
		connectionIP = getNewPaxosAddressBank(conn.RemoteAddr().String())
		conn1, _ = net.Dial("tcp", connectionIP)
	}
}

func waitForResponseBank() {
	service := "0.0.0.0:1337"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	Utils.CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	for {
		Utils.CheckError(err)
		conn, _ := listener.Accept()
		go holdClientConnectionBank(conn)
	}
}

func getNewPaxosAddressBank(failedAddress string) string {
	var newAddress string
R:
	for _, v := range serverListBank {
		if v.IP != Utils.GetIp(failedAddress) {
			newAddress = v.IP
			fmt.Println("New address:", newAddress)
			break R
		}
	}
	return newAddress
}

func holdClientConnectionBank(conn net.Conn) {
	var connectionOK = true
	decoder := gob.NewDecoder(conn)
	for connectionOK == true {
		var message interface{}
		err := decoder.Decode(&message)
		switch message.(type) {
		case msg.ClientResponseMessage:
			if err != nil {
				fmt.Println("Closed?")
				connectionOK = false
			}
			if message != nil {
				var clientMsg msg.ClientResponseMessage
				clientMsg = message.(msg.ClientResponseMessage)
				lastConfirmedValueBank = clientMsg.Value
				lastConfirmedMsgNumberBank = clientMsg.MsgNumber
				var stringMessage = fmt.Sprintf("Message from bank: %s ", clientMsg.Value)
				fmt.Println("---------------------------------------------------")
				fmt.Println(stringMessage)
				fmt.Println("---------------------------------------------------")
				if sendAllBank == false {
					sentChanBank <- 1
				}
			} else {
				fmt.Println("Message is empty stupid!")
				connectionOK = false
			}
		case msg.ClientResponseNodes:
			var clientMsg msg.ClientResponseNodes
			clientMsg = message.(msg.ClientResponseNodes)
			serverListBank = clientMsg.List
		}
	}
	fmt.Println("Paxos closed connection, no more to share")
}

func GetServersBank(conn net.Conn, myIP string) {
	for {
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(5000 * time.Millisecond)
			timeout <- true
		}()
		select {
		case <-timeout:
			fmt.Println("Sending request")
			encoder := gob.NewEncoder(conn1)
			var sendMsg = msg.ClientRequestNodes{RemoteAddress: myIP + ":1337"}
			var message interface{}
			message = sendMsg
			err := encoder.Encode(&message)
			if err != nil {
				fmt.Println("Finding new address")
				connectionIP = getNewPaxosAddressBank(conn.RemoteAddr().String())
				conn1, _ = net.Dial("tcp", connectionIP)
			}
		}
	}
}

func GetIPBank() string {
	name, _ := os.Hostname()
	addr, _ := net.LookupHost(name)
	UDPAddr, _ := net.ResolveUDPAddr("udp4", addr[0]+":1888")
	ip := UDPAddr.IP.String()
	return ip
}
