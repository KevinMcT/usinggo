package msgsServer

import (
	"encoding/gob"
	"fmt"
	"lab2/helpers"
	"lab2/messages"
	"net"
	"os"
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

	service := "0.0.0.0:" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	CheckError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	CheckError(err)

	conn, err := listener.Accept()

	for {
		if err != nil {
			continue
		}

		decoder := gob.NewDecoder(conn)
		demarshal(decoder, inChan)
	}
	conn.Close()
}

func demarshal(dec *gob.Decoder, inChan chan interface{}) {
	var msg interface{}
	err := dec.Decode(&msg)
	if err != nil {
		fmt.Println(err)
		fmt.Println("is this error?")
	}
	inChan <- msg
}

func handleMsgs(inchan chan messages.StrMsg, quitMsgs chan int) {
	senders := ""
	for {
		msg, ok := <-inchan
		if ok {
			senders = senders + msg.Sender + "\nWith content: " + msg.Content + "\n"
		} else {
			if senders != "" {
				fmt.Println("The senders: \n" + senders)
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
			errors = errors + msg.Sender + "\nWith content: " + msg.Error + "\n"
		} else {
			if errors != "" {
				fmt.Println("The errors from: \n" + errors)
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
