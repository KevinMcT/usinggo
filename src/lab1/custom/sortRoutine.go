package custom

import (
	"lab1/messages"
)

func Sort(inChan chan interface{}, msgChan chan messages.StrMsg, errChan chan messages.ErrMsg) {
	//receive messages from inChan
	//forward messages on msgChan or errChan according to its type. 
	//  (use switch x.(type) { case:...})
	//when all was send, close the channels
	for msg, ok := <-inChan; ok; msg, ok = <-inChan {

		switch msg.(type) {
		case messages.StrMsg:
			msgChan <- msg.(messages.StrMsg)
		case messages.ErrMsg:
			errChan <- msg.(messages.ErrMsg)
		}
	}
	close(errChan)
	close(msgChan)
}
