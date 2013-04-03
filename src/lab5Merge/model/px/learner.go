package px

import (
	"encoding/gob"
	"fmt"
	"lab5Merge/model/RoundVar"
	"lab5Merge/model/net/msg"
	"lab5Merge/model/net/tcp"
	//"net"
	"strings"
	"time"
)

var (
	learns        int
	value         string
	learnList     = make([]msg.Learn, 0)
	r             int
	msgNr         int
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
		learn := <-msg.LearnChan
		if waiting == false {
			waitLearnChan <- "wait"
		}
		learnMsg := learn.Message.(msg.Learn)
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
				r = v.ROUND
				msgNr = v.MSGNUMBER
				learns = learns + 1
			} else if strings.EqualFold(v.VALUE, value) == true {
				learns = learns + 1
			} else {
				value = v.VALUE
				r = v.ROUND
				msgNr = v.MSGNUMBER
				learns = learns + 1
			}
		}
		if learns > (len(RoundVar.GetRound().List) / 2) {
			var stringMessage = fmt.Sprintf("Learnt value %s round:%d messageNumber:%d ", value, r, msgNr)
			learns = 0
			value = "-1"
			fmt.Println(stringMessage)
			if leader.IP == self.IP {
				sendAddress := RoundVar.GetRound().RespondClient + ":1337"
				sendConn := tcp.Dial(sendAddress)
				if sendConn != nil {
					encoder := gob.NewEncoder(sendConn)
					var prepare = msg.ClientResponseMessage{Content: stringMessage}
					var message interface{}
					message = prepare
					encoder.Encode(&message)
					tcp.Close(sendConn)
				}
			}

		} else {
			fmt.Println("Did not receive not enough learns, not learning anything!")
		}
	}
}
