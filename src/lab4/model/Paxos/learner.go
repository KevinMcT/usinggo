package Paxos

import (
	//"encoding/gob"
	"fmt"
	"lab4/model/Network/message"
	//"net"
)

func Learner() {
	fmt.Println("Learner up and waiting ...")
	go receivedLearn()
}

func receivedLearn() {
	for {
		value := <-message.LearnChan
		fmt.Println("Im going to learn the value:")
		fmt.Println(value)
	}
}
