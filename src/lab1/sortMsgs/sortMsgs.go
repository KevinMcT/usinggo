//use: sortMsgs [file]
package main

import (
	"encoding/gob"
	"fmt"
	"lab1/custom"
	"lab1/messages"
	"os"
)

var (
	inChan   = make(chan interface{}, 5)     //channel from demarshaler to sort	
	msgChan  = make(chan messages.StrMsg, 5) //message channel
	errChan  = make(chan messages.ErrMsg, 5) //error channel
	quitMsgs = make(chan int)
	quitErr  = make(chan int)
)

func main() {
	go handleMsgs(msgChan, quitMsgs)
	go handleErrors(errChan, quitErr)
	go custom.Sort(inChan, msgChan, errChan)
	demarshal(os.Args[1], inChan)
	<-quitMsgs
	<-quitErr
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
	close(inChan)
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
