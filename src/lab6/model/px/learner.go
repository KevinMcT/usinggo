package px

import (
	"fmt"
	"lab6/Utils"
	"lab6/model/RoundVar"
	"lab6/model/SlotList"
	"lab6/model/net/msg"
	"lab6/model/net/tcp"
	"time"
)

var (
	learns              int
	value               interface{}
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
	value = nil
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
		fmt.Println(learnList)
		fmt.Println("Init value: ", value)
		for _, v := range learnList {
			if v.ROUND >= r {
				if value == nil {
					value = v.VALUE
					fmt.Println("new Value: ", value)
					r = v.ROUND
					msgNr = v.MSGNUMBER
					learns = learns + 1
				} else if Utils.Equals(value, v.VALUE) && r == v.ROUND && msgNr == v.MSGNUMBER {
					learns = learns + 1
				}
			}
		}
		if learns > (len(RoundVar.GetRound().List) / 2) {
			var stringMessage = fmt.Sprintf("Learnt value %s round:%d messageNumber:%d ", value, r, msgNr)
			learns = 0
			lastLearntMsgNumber = msgNr
			learnList = make([]msg.Learn, 0)
			fmt.Println("--", stringMessage, "--")
			var addedSlot = slots.Add(value, msgNr)
			if addedSlot {
				fmt.Println("Wrote to slot: ", value)
			}
			if RoundVar.GetRound().CurrentLeader.IP == self.IP {
				sendAddress := RoundVar.GetRound().RespondClient + ":1337"
				var prepare = msg.ClientResponseMessage{Value: handleBankRequest(value), Round: r, MsgNumber: msgNr}
				var message interface{}
				message = prepare
				tcp.SendPaxosMessage(sendAddress, message)
			}
			value = nil
		} else {
			fmt.Println("--Did not receive not enough learns, not learning anything!--")
		}
	}
}

func handleBankRequest(req interface{}) string {
	switch req.(type) {
	case msg.Deposit:
		dep := req.(msg.Deposit)
		if dep.Amount < 0 {
			return fmt.Sprintf("Cannot deposit negative cash. This does not exist, even in NBB")
		}
		oldBalance := bankAccounts[dep.AccountNumber]
		newBalance := oldBalance + dep.Amount
		bankAccounts[dep.AccountNumber] = newBalance
		return fmt.Sprintf("Deposited %d into %s, new balance %d", dep.Amount, dep.AccountNumber, newBalance)
	case msg.Withdraw:
		rem := req.(msg.Withdraw)
		oldBalance := bankAccounts[rem.AccountNumber]
		newBalance := oldBalance - rem.Amount
		if newBalance < 0 {
			return fmt.Sprintf("Insufficiant funds in Account %s", rem.Amount)
		}
		if rem.Amount < 0 {
			return fmt.Sprintf("Cannot withdraw negative cash. This does not exist, even in NBB")
		}
		bankAccounts[rem.AccountNumber] = newBalance
		return fmt.Sprintf("Withdrew %d out of %s, new balance %d", rem.Amount, rem.AccountNumber, newBalance)
	case msg.Transfer:
		tran := req.(msg.Transfer)
		//From
		oldBalance := bankAccounts[tran.FromAccount]
		if tran.Amount < 0 {
			return fmt.Sprintf("Cannot transfer negative cash. This does not exist, even in NBB")
		}
		newBalance := oldBalance - tran.Amount
		if newBalance < 0 {
			return fmt.Sprintf("Insufficiant funds in transfer Account %s", tran.FromAccount)
		}
		bankAccounts[tran.FromAccount] = newBalance
		//To
		oldBalance = bankAccounts[tran.ToAccount]
		newBalance = oldBalance + tran.Amount
		bankAccounts[tran.ToAccount] = newBalance

		return fmt.Sprintf("Transferd %d from %s to %s", tran.Amount, tran.FromAccount, tran.ToAccount)
	case msg.Balance:
		bal := req.(msg.Balance)
		balance := bankAccounts[bal.AccountNumber]
		return fmt.Sprintf("Balance on account %s: %d", bal.AccountNumber, balance)
	}
	return ""
}
