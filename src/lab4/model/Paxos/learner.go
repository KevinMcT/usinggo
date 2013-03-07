package Paxos

import (
	//"encoding/gob"
	"fmt"
	"lab4/model/Network/message"
	//"net"
	"strings"
)

var (
	learns int
	value  string
)

func Learner() {
	fmt.Println("Learner up and waiting ...")
	learns = -1
	value = ""
	go receivedLearn()
}

func receivedLearn() {
	for {
		learn := <-message.LearnChan
		learnMsg := learn.Message.(message.Learn)
		if learns == -1 {
			value = learnMsg.VALUE
			learns = 1
		} else {
			if strings.EqualFold(learnMsg.VALUE, value) == true {
				learns = learns + 1
			}
			if learns > ((len(nodeList) / 2) + 1) {
				fmt.Println("Learnt value ", value)
				learns = -1
				value = ""
			}
		}
	}
}
