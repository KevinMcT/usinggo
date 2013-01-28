package helpers

import "lab2/messages"

func Sort(inChan chan interface{}, msgChan chan messages.StrMsg, errChan chan messages.ErrMsg) {
	for {
		m, ok := <-inChan
		if ok {
			switch m.(type) {
			case messages.StrMsg:
				msgChan <- m.(messages.StrMsg)
			case messages.ErrMsg:
				errChan <- m.(messages.ErrMsg)
			case messages.NT:
				//DO NOTHING! DO NOT TOUCH!! DISCARDED PIIIINGZZZZ!
			}
		} else {
			break
		}
	}

	close(errChan)
	close(msgChan)
}

//func Sort(inChan chan interface{}, msgChan chan messages.StrMsg, errChan chan messages.ErrMsg) {
//	//receive messages from inChan
//	//forward messages on msgChan or errChan according to its type. 
//	//  (use switch x.(type) { case:...})
//	//when all was send, close the channels
//	for msg, ok := <-inChan; ok; msg, ok = <-inChan {
//		switch msg.(type) {
//		case messages.StrMsg:
//			msgChan <- msg.(messages.StrMsg)
//		case messages.ErrMsg:
//			errChan <- msg.(messages.ErrMsg)
//		}
//	}
//	close(errChan)
//	close(msgChan)
//}
