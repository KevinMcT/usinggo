package Paxos

import (
	//"encoding/gob"
	"fmt"
	"lab4/model/Network/message"
	"lab4/model/RoundVar"
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
		learnList = append(learnList, learnMsg)
	}
}

/*
After receiving the first learn we wait for the rest. If
we have received more then halv of all node learns, we learn
the value
*/
func waitForLearns() {
	for {
		<-waitLearnChan
		waiting = true
		time.Sleep(2 * time.Second)
		waiting = false
		for _, v := range learnList {
			if strings.EqualFold(value, "-1") == true {
				value = v.VALUE
				learns = learns + 1
			} else if strings.EqualFold(v.VALUE, value) == true {
				learns = learns + 1
			}
		}
		if learns > (len(RoundVar.GetRound().List) / 2) {
			fmt.Println("Learnt value ", value)
			learns = 0
			value = "-1"
			learnList = make([]message.Learn, 0)
		}
	}
}
