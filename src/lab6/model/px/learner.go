package px

import (
	"fmt"
	"lab6/model/RoundVar"
	"lab6/model/SlotList"
	"lab6/model/net/msg"
	"lab6/model/net/tcp"
	"strings"
	"time"
)

var (
	learns              int
	value               string
	learnList           = make([]msg.Learn, 0)
	r                   int
	msgNr               int
	slots               *SlotList.Slots
	bankAccounts        = make(map[string]int, 0)
	waitLearnChan       = make(chan string, 1)
	lastLearntMsgNumber int
)

func Learner(slotList *SlotList.Slots) {
	fmt.Println("--Learner up and waiting ...")
	learns = 0
	value = "-1"
	lastLearntMsgNumber = -1
	slots = slotList
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
		time.Sleep(15 * time.Millisecond)
		waiting = false
		for _, v := range learnList {
			if v.ROUND >= r {
				if strings.EqualFold(value, "-1") {
					value = v.VALUE
					r = v.ROUND
					msgNr = v.MSGNUMBER
					learns = learns + 1
				} else if strings.EqualFold(v.VALUE, value) && r == v.ROUND && msgNr == v.MSGNUMBER {
					learns = learns + 1
				}
				/*var addedSlot = slots.Add(v, v.MSGNUMBER-1)
				if addedSlot {
					fmt.Println("Wrote to slot: ", v)
				}*/
			}
		}
		if learns > (len(RoundVar.GetRound().List) / 2) {
			var stringMessage = fmt.Sprintf("Learnt value %s round:%d messageNumber:%d ", value, r, msgNr)
			learns = 0
			lastLearntMsgNumber = msgNr
			learnList = make([]msg.Learn, 0)
			fmt.Println("--", stringMessage, "--")
			if RoundVar.GetRound().CurrentLeader.IP == self.IP {
				sendAddress := RoundVar.GetRound().RespondClient + ":1337"
				var prepare = msg.ClientResponseMessage{Value: value, Round: r, MsgNumber: msgNr}
				var message interface{}
				message = prepare
				tcp.SendPaxosMessage(sendAddress, message)
			}
			value = "-1"
		} else {
			fmt.Println("--Did not receive not enough learns, not learning anything!--")
		}
	}
}
