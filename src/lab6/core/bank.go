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
	fmt.Println(ip)
	sendAllBank = true
	for {
		fmt.Println("Connecting to Paxos replica")
		service := ip + ":1337"
		paxosAddressBank = service
		conn, err := net.Dial("tcp", paxosAddressBank)
		fmt.Println("Waiting for servers... This might take up to 5 seconds, but you can still send a message")
		lastConfirmedValueBank = ""
		lastConfirmedMsgNumberBank = -1
		go GetServersBank(conn, GetIPBank())
		if err == nil {
			fmt.Println("*********************************************************************")
			fmt.Println("*          Welcome to the National Bank of Bullshit                 *")
			fmt.Println("* Choose the function you would like to test using numbers 1 - 2    *")
			fmt.Println("* 1 - Deposit                                                       *")
			fmt.Println("* 2 - Withdraw                                                     *")
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
				fmt.Scanf("%d", &accFrm)
				fmt.Println("-- Enter amount you wish to withdraw --")
				fmt.Scanf("%d", &amt)
				bankMessage = msg.Withdraw{AccountNumber: accFrm, Amount: amt}
				break
			case 3:
				fmt.Println("-- Enter your account number --")
				fmt.Scanf("%d", &accFrm)
				fmt.Println("-- Enter amount you wish to transfer --")
				fmt.Scanf("%d", &amt)
				fmt.Println("-- Enter account to transfer to --")
				fmt.Scanf("%d", accTo)
				bankMessage = msg.Transfer{FromAccount: accFrm, ToAccount: accTo, Amount: amt}
				break
			case 4:
				fmt.Println("-- Enter your account number --")
				fmt.Scanf("%d", &accFrm)
				bankMessage = msg.Balance{AccountNumber: accFrm}
				break
			case 0:
				os.Exit(0)
			}
			sendToPaxosBank(bankMessage, conn, 0, 1)
		} else {
			fmt.Println("--Seems like the node you are trying to connect is gone down or does not exist. Please try another address--")
		}
	}
}

func sendToPaxosBank(st interface{}, conn net.Conn, start int, end int) {
	var paxosConn = conn
	var allOk = true
L:
	for i := start; i < end; i++ {
		encoder := gob.NewEncoder(paxosConn)
		var sendMsg = msg.ClientRequestMessage{Content: st}
		var message interface{}
		message = sendMsg
		var err = encoder.Encode(&message)
		if err != nil {
			var address = getNewPaxosAddressBank(Utils.GetIp(conn.RemoteAddr().String()))
			address = address + ":1337"
			paxosAddressBank = address
			time.Sleep(1000 * time.Millisecond)
			newConn, _ := net.Dial("tcp", paxosAddressBank)
			paxosConn = newConn
			allOk = false
			break L
		}
		if sendAllBank == false {
			timeout := make(chan bool, 1)
			go func() {
				time.Sleep(5000 * time.Millisecond)
				timeout <- true
			}()
			select {
			case <-sentChanBank:
				//Don`t do anything here				
			case <-timeout:
				fmt.Println("--No reply on message from connection, finding a new one!--")
				var address = getNewPaxosAddressBank(Utils.GetIp(conn.RemoteAddr().String()))
				address = address + ":1337"
				paxosAddressBank = address
				time.Sleep(1000 * time.Millisecond)
				newConn, _ := net.Dial("tcp", paxosAddressBank)
				paxosConn = newConn
				allOk = false
				break L
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	if allOk == false {
		fmt.Println("Starting a new send from last learnt value!")
		fmt.Println("Last confirmed message: ", lastConfirmedMsgNumberBank)
		sendToPaxosBank(st, paxosConn, lastConfirmedMsgNumberBank+1, end)
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
	for _, v := range serverListBank {
		if v.IP != failedAddress {
			newAddress = v.IP
			break
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
				connectionOK = false
			}
			if message != nil {
				var clientMsg msg.ClientResponseMessage
				clientMsg = message.(msg.ClientResponseMessage)
				lastConfirmedValueBank = clientMsg.Value
				lastConfirmedMsgNumberBank = clientMsg.MsgNumber
				var stringMessage = fmt.Sprintf("Learnt value %s round:%d messageNumber:%d", clientMsg.Value, clientMsg.Round, clientMsg.MsgNumber)
				fmt.Println("---------------------------------------------------")
				fmt.Println(stringMessage)
				fmt.Println("---------------------------------------------------")
				if sendAllBank == false {
					sentChanBank <- 1
				}
			} else {
				fmt.Println("Message is empty stupid!")
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
			encoder := gob.NewEncoder(conn)
			var sendMsg = msg.ClientRequestNodes{RemoteAddress: myIP + ":1337"}
			var message interface{}
			message = sendMsg
			encoder.Encode(&message)
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
