//use: sortMsgs [file]
package main

import (
	"encoding/gob"
	"fmt"
	"messages"
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
	file, _ := os.Open(filePath)
	inputMessages := make([]interface{}, 10)
	decoder := gob.NewDecoder(file)
	for {
		inputMessages = append(inputMessages, decoder.Decode(i))
	}
}

func sort(inChan chan interface{}, msgChan chan messages.StrMsg, errChan chan messages.ErrMsg) {
	//receive messages from inChan
	//forward messages on msgChan or errChan according to its type. 
	//  (use switch x.(type) { case:...})
	//when all was send, close the channels
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
	}
}
