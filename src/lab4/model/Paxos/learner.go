package Paxos

import (
	//"encoding/gob"
	"fmt"
	"lab4/model/Network/message"
	//"net"
	"strings"
	"time"
)

var (
	learns        int
	value         string
	learnList     = make([]message.Learn, 0)
	waitLearnChan = make(chan string, 1)
)

func Learner() {
	fmt.Println("Learner up and waiting ...")
	learns = 0
	value = "-1"
	go receivedLearn()
	go waitForLearns()
}

func receivedLearn() {
	for {
		learn := <-message.LearnChan
		if waiting == false {
			waitLearnChan <- "wait"
		}
		learnMsg := learn.Message.(message.Learn)
		fmt.Println("Received learn: ", learnMsg)
		learnList = append(learnList, learnMsg)
		/*if learns == -1 {
			value = learnMsg.VALUE
			fmt.Println("Hoping for more: ", value)
			learns = 1
		} else {
			fmt.Println("stored learn: ", value)
			fmt.Println("Received learn: ")
			if strings.EqualFold(learnMsg.VALUE, value) == true {
				learns = learns + 1
			}
			if learns > (len(nodeList) / 2) {
				fmt.Println("Learnt value ", value)
				learns = -1
				value = ""
			} else {
				fmt.Println("Learns: ", learns)
				fmt.Println("Req. learns: ", (len(nodeList) / 2))
			}
		}*/
	}
}

func waitForLearns() {
	for {
		<-waitLearnChan
		waiting = true
		time.Sleep(2 * time.Second)
		waiting = false
		fmt.Println(learnList)
		for _, v := range learnList {
			if strings.EqualFold(value, "-1") == true {
				value = v.VALUE
				learns = learns + 1
			} else if strings.EqualFold(v.VALUE, value) == true {
				learns = learns + 1
			}
		}
		if learns > (len(nodeList) / 2) {
			fmt.Println("Learnt value ", value)
			learns = 0
			value = "-1"
			learnList = make([]message.Learn, 0)
		}
	}
}
