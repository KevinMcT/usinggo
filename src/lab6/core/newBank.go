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
	nodeServersList   []node.T_Node
	udpListenChan1    = make(chan string, 0)
	newConnRequired   = make(chan bool, 0)
	firstConnection   = make(chan int, 0)
	newServerDetected = make(chan int, 0)
	conn              net.Conn
	err               error
	ip                string
	pMessage          interface{}
	first             bool
)

func NewBank() {
	first = true
	go udp.Listen(udpListenChan1)
	go ListenForConnections()
	go HandleNewConnection()
	udp.SendLocator(GetIPBank())
	ip = <-udpListenChan1
	conn, err = CreateDialUp(ip)
	go GetServerNodes(GetMyIP())
	Runnable()
}

func HandleNewConnection() {
	for {
		<-newConnRequired
		for i, v := range nodeServersList {
			if ip == Utils.GetIp(v.IP) {
				nodeServersList[i].SUSPECTED = true
			}
		}
		ip = GetNewAddress(ip)
		conn, err = CreateDialUp(ip)
		fmt.Println("New connection to: ", ip, "with conn:", conn, "established")
		newServerDetected <- 1
	}
}

func Runnable() {
	for {
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(50 * time.Millisecond)
			timeout <- true
		}()
		select {
		case <-timeout:
			fmt.Println("*********************************************************************")
			fmt.Println("*          Welcome to the National Bank of Bullshit                 *")
			fmt.Println("* Choose the function you would like to test using numbers 1 - 4    *")
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
				pMessage = msg.Deposit{AccountNumber: accFrm, Amount: amt}
				break
			case 2:
				fmt.Println("-- Enter your account number --")
				fmt.Scanf("%s", &accFrm)
				fmt.Println("-- Enter amount you wish to withdraw --")
				fmt.Scanf("%d", &amt)
				pMessage = msg.Withdraw{AccountNumber: accFrm, Amount: amt}
				break
			case 3:
				fmt.Println("-- Enter your account number --")
				fmt.Scanf("%s", &accFrm)
				fmt.Println("-- Enter amount you wish to transfer --")
				fmt.Scanf("%d", &amt)
				fmt.Println("-- Enter account to transfer to --")
				fmt.Scanf("%s", &accTo)
				pMessage = msg.Transfer{FromAccount: accFrm, ToAccount: accTo, Amount: amt}
				break
			case 4:
				fmt.Println("-- Enter your account number --")
				fmt.Scanf("%s", &accFrm)
				pMessage = msg.Balance{AccountNumber: accFrm}
				break
			case 0:
				os.Exit(0)
			}
			SendMessage(pMessage)
		}
	}
}

func GetServerNodes(myIP string) {
	for {
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(5000 * time.Millisecond)
			timeout <- true
		}()
		select {
		case <-timeout:
			var sendMsg = msg.ClientRequestNodes{RemoteAddress: myIP + ":1337"}
			var message interface{}
			message = sendMsg
			enc := gob.NewEncoder(conn)
			erro := enc.Encode(&message)
			if erro != nil {
				fmt.Println("AUTO: New server required")
				newConnRequired <- true
			}
		}
	}
}

func SendMessage(mesg interface{}) {
	var sendMsg = msg.ClientRequestMessage{Content: mesg}
	var message interface{}
	message = sendMsg
	enc := gob.NewEncoder(conn)
	var err = enc.Encode(&message)
	if err != nil {
		fmt.Println("MSG: New server required")
		newConnRequired <- true
		<-newServerDetected
		SendMessage(mesg)
	}
}

func CreateDialUp(ip string) (net.Conn, error) {
	return net.Dial("tcp", ip+":1337")
}

func GetMyIP() string {
	name, _ := os.Hostname()
	addr, _ := net.LookupHost(name)
	UDPAddr, _ := net.ResolveUDPAddr("udp4", addr[0]+":1888")
	ip := UDPAddr.IP.String()
	return ip
}

func GetNewAddress(failedAddress string) string {
	var newAddress string
R:
	for _, v := range nodeServersList {
		if v.IP != Utils.GetIp(failedAddress) && v.SUSPECTED != true {
			newAddress = v.IP
			break R
		}
	}
	return newAddress
}

func ListenForConnections() {
	service := "0.0.0.0:1337"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	Utils.CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	for {
		Utils.CheckError(err)
		conn, _ := listener.Accept()
		go KeepConn(conn)
	}
}

func KeepConn(conn net.Conn) {
	var connectionOK = true
	decoder := gob.NewDecoder(conn)
	for connectionOK == true {
		var message interface{}
		erro := decoder.Decode(&message)
		switch message.(type) {
		case msg.ClientResponseMessage:
			if erro != nil {
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
			nodeServersList = clientMsg.List
			if first {
				fmt.Println("Servers is downloaded")
				first = false
			}
		}
	}
	fmt.Println("Paxos closed connection, no more to share")
}
