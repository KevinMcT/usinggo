package msgsServer

import (
	"encoding/gob"
	"fmt"
	"lab2/helpers"
	"lab2/messages"
	"net"
	"os"
	"time"
)

var (
	inChan   = make(chan interface{}, 5)     //channel from demarshaler to sort	
	msgChan  = make(chan messages.StrMsg, 5) //message channel
	errChan  = make(chan messages.ErrMsg, 5) //error channel
	quitMsgs = make(chan int)
	quitErr  = make(chan int)
)

func MsgsServer(port string) {
	go handleMsgs(msgChan, quitMsgs)
	go handleErrors(errChan, quitErr)
	go helpers.Sort(inChan, msgChan, errChan)
	demarshal(port, inChan)
	<-quitMsgs
	<-quitErr
}

func demarshal(port string, inChan chan interface{}) {
	service := "0.0.0.0:" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	CheckError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	CheckError(err)

	conn, err := listener.Accept()

	for {
		conn.SetDeadline(time.Now().Add(30 * time.Second))
		if err != nil {
			break
		}
		decoder := gob.NewDecoder(conn)
		var msg interface{}
		err = decoder.Decode(&msg)
		if err != nil {
			fmt.Println(err)
		}
		if msg != nil {
			conn.SetDeadline(time.Now().Add(30 * time.Second))
		}
		inChan <- msg
	}
	close(inChan)
}

func handleMsgs(inchan chan messages.StrMsg, quitMsgs chan int) {
	senders := ""
	for {
		msg, ok := <-inchan
		if ok {
			senders = senders + msg.Sender + ", "
			fmt.Println("Message recived from: ", msg.Sender, ":", msg.Content)
		} else {
			if senders != "" {
				fmt.Println("The senders: " + senders)
			}
			break
		}
	}
	quitMsgs <- 1
}

func handleErrors(inchan chan messages.ErrMsg, quitErr chan int) {
	errors := ""
	for {
		msg, ok := <-inchan
		if ok {
			errors = errors + msg.Error + ", "
			fmt.Println("Message recived from: ", msg.Sender, ":", msg.Error)
		} else {
			if errors != "" {
				fmt.Println("The errors from: " + errors)
			}
			break
		}
	}
	quitErr <- 1
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
