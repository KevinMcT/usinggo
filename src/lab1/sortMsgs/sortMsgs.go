//use: sortMsgs [file]
package main

import (
	"encoding/gob"
	"fmt"
	"lab1/messages"
	"os"
	"time"
)

var (
	inChan  = make(chan interface{}, 5)     //channel from demarshaler to sort	
	msgChan = make(chan messages.StrMsg, 5) //message channel
	errChan = make(chan messages.ErrMsg, 5) //error channel
)

func main() {

	go handleMsgs(msgChan)
	go handleErrors(errChan)
	go sort(inChan, msgChan, errChan)
	demarshal(os.Args[1], inChan)
	time.Sleep(5 * time.Second)
}

func demarshal(filePath string, inChan chan interface{}) {
	//Retrieve the slice of messages from the file (type []interface{})
	//send each element of the slice on inChan
	//when all was send, close the channel.
	input, _ := os.Open(filePath)
	dec := gob.NewDecoder(input)
	t2 := make([]interface{}, 5)
	err := dec.Decode(&t2)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range t2 {
		inChan <- v
	}
}

func sort(inChan chan interface{}, msgChan chan messages.StrMsg, errChan chan messages.ErrMsg) {
	//receive messages from inChan
	//forward messages on msgChan or errChan according to its type. 
	//  (use switch x.(type) { case:...})
	//when all was send, close the channels
	for {
		msg, _ := <-inChan
		switch msg.(type) {
		case messages.StrMsg:
			msgChan <- msg.(messages.StrMsg)
		case messages.ErrMsg:
			errChan <- msg.(messages.ErrMsg)
		}
	}
}

func handleMsgs(inchan chan messages.StrMsg) {
	senders := ""
	for {
		msg, ok := <-inchan
		if ok {
			senders = senders + msg.Sender + ", "
		} else {
			if senders != "" {
				fmt.Println("The senders: " + senders)
			}
			break
		}
		fmt.Println("Senderen = " + senders)
	}
}

func handleErrors(inchan chan messages.ErrMsg) {
	errors := ""
	for {
		msg, ok := <-inchan
		if ok {
			errors = errors + msg.Error + ", "
		} else {
			if errors != "" {
				fmt.Println("The errors: " + errors)
			}
			break
		}
		fmt.Println("Feilmld = " + errors)
	}
}
