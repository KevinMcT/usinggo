package fd

import (
	"btb/controller/node"
	"btb/model/net/msg"
	"btb/model/net/tcp"
	"fmt"
	"time"
)

var (
	delay int
)

func Detect(me node.T_Node, lead node.T_Node, newLead chan node.T_Node, newNodeChan chan node.T_Node, suspectedChan chan node.T_Node, restoreChan chan node.T_Node, tcpRequestChan chan msg.HARTBEATREQUEST, tcpResponseChan chan msg.HARTBEATRESPONSE, startList []node.T_Node, endList chan []node.T_Node) {
	var ticker = time.NewTicker(50 * time.Millisecond)
	var nodeList = startList
	var newNode node.T_Node
	var newL node.T_Node
	go GetResponse(tcpResponseChan)
	for {
		time.Sleep(1 * time.Millisecond)
		if me.IP == lead.IP {
			for i, v := range nodeList {
				if v.IP != me.IP && v.LEAD != true {
					err := tcp.Send(v.IP, msg.HARTBEATREQUEST{IP: lead.IP})
					time.Sleep(5 * time.Millisecond)
					if err != nil && nodeList[i].SUSPECTED != true {
						fmt.Println("FD: suspect...")
						nodeList[i].SUSPECTED = true
						suspectedChan <- nodeList[i]
					}
					if err == nil && nodeList[i].SUSPECTED == true {
						fmt.Println("FD: restore...")
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
			nodeList = AppendIfMissing(nodeList, newNode)
		case <-timeout:
		case newL = <-newLead:
			fmt.Println("New leader detected in FD:", newL)
			lead = newL
		case list := <-endList:
			nodeList = list
		}
		<-ticker.C
	}
}
func AppendIfMissing(slice []node.T_Node, i node.T_Node) []node.T_Node {
	for j, ele := range slice {
		if ele.IP == i.IP {
			slice[j].TIME = i.TIME
			return slice
		}
	}
	return append(slice, i)
}

func GetResponse(tcpResponseChan chan msg.HARTBEATRESPONSE) {
	var i int
	for {
		<-tcpResponseChan
		i = i + 1
	}
}
