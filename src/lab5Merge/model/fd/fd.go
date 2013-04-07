package fd

import (
	"fmt"
	"lab5Merge/controller/node"
	"lab5Merge/model/net/msg"
	"lab5Merge/model/net/tcp"
	"time"
)

var (
	delay int
)

/*
	Sends hartbeats to all the nodes in the nodelist. If the response takes to long to arrive the
	TCP class will throw an error, which is caught here. This will result in the node being marked as suspected.
*/

/*
	Public function
	Should be run as go-routine
	me 					node.T_Node 				This node
	lead 				node.T_Node 				System leader
	newLead				chan node.T_Node			If new leader, this channel handles it
	newNodeChan			chan node.T_Node			When a new node is detected it should be added on this channel
	suspectedChan		chan node.T_Node			If a node is suspected, it will be sent on this channel
	restoreChan			chan node.T_Node			If a node is restored, it will be sent on this channel
	tcpRequestChan		chan msg.HARTBEATREQUEST	TCP class will send the request on this channel. The slaves reads it and sends a tcp response
	tcpResponseChan		chan msg.HARTBEATRESPONSE	Is not used to other than empty the buffer of the tcp recieve.		
	startList			[]node.T_Node				The list the system has on startup
	endList				chan []node.T_Node			Upon update of the entire list, it will be sent on this channel.
*/

func Detect(me node.T_Node, lead node.T_Node, newLead chan node.T_Node, newNodeChan chan node.T_Node, suspectedChan chan node.T_Node, restoreChan chan node.T_Node, tcpRequestChan chan msg.HARTBEATREQUEST, tcpResponseChan chan msg.HARTBEATRESPONSE, startList []node.T_Node, endList chan []node.T_Node) {
	var ticker = time.NewTicker(50 * time.Millisecond)
	var nodeList = startList
	var newNode node.T_Node
	var newL node.T_Node
	go getResponse(tcpResponseChan)
	for {
		time.Sleep(1 * time.Millisecond)
		if me.IP == lead.IP {
			for i, v := range nodeList {
				if v.IP != me.IP && v.LEAD != true {
					err := tcp.Send(v.IP, msg.HARTBEATREQUEST{IP: lead.IP})
					time.Sleep(5 * time.Millisecond)
					if err != nil && nodeList[i].SUSPECTED != true {
						nodeList[i].SUSPECTED = true
						suspectedChan <- nodeList[i]
					}
					if err == nil && nodeList[i].SUSPECTED == true {
						nodeList[i].SUSPECTED = false
					}
					time.Sleep(1 * time.Millisecond)
				}
			}
		}
		if me.IP != lead.IP {
			timeout := make(chan bool, 1)
			go func() {
				time.Sleep(5 * time.Millisecond)
				timeout <- true
			}()
			select {
			case <-tcpRequestChan:
				tcp.Send(lead.IP, msg.HARTBEATRESPONSE{IP: me.IP})
				time.Sleep(5 * time.Millisecond)
			case <-timeout:
			}
		}

		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(50 * time.Millisecond)
			timeout <- true
		}()
		select {
		case newNode = <-newNodeChan:
			nodeList = appendIfMissing(nodeList, newNode)
		case <-timeout:
		case newL = <-newLead:
			fmt.Println("--New leader detected in FD:", newL, "--")
			lead = newL
		case list := <-endList:
			nodeList = list
		}
		<-ticker.C
	}
}

/*
	Private function
	slice 		[]node.T_Node 	List to be checked
	i 			node.T_Node 	Element to be added
	returns 	[]node.T_Node 	The new list
*/
func appendIfMissing(slice []node.T_Node, i node.T_Node) []node.T_Node {
	for j, ele := range slice {
		if ele.IP == i.IP {
			slice[j].TIME = i.TIME
			return slice
		}
	}
	return append(slice, i)
}

/*
	Private function
	tcpResponseChan chan msg.HARTBEATRESPONSE Channel to empty
*/
func getResponse(tcpResponseChan chan msg.HARTBEATRESPONSE) {
	for {
		<-tcpResponseChan
	}
}
