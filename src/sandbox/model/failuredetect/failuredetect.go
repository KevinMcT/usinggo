package failuredetect

import (
	"fmt"
	"sandbox/controller/machine"
	"sandbox/controller/node"
	"sandbox/model/network/message"
	"sandbox/model/network/tcp"
	"time"
)

var (
	delay int
)

func Detect(me machine.T_Machine, lead node.T_Node, newLead chan node.T_Node, newNodeChan chan node.T_Node, suspectedChan chan node.T_Node, restoreChan chan node.T_Node, tcpRequestChan chan message.HARTBEATREQUEST, tcpResponseChan chan message.HARTBEATRESPONSE, startList []node.T_Node, endList chan []node.T_Node) {
	var ticker = time.NewTicker(200 * time.Millisecond)
	var nodeList = startList
	var newNode node.T_Node
	var newL node.T_Node
	for {
		time.Sleep(1 * time.Millisecond)
		if me.IP == lead.IP {
			for i, v := range nodeList {
				if v.IP != me.IP && v.LEAD != true {
					err := tcp.Send(v.IP, message.HARTBEATREQUEST{IP: lead.IP})
					fmt.Println("SENT REQUEST TO", v.IP)
					if err != nil && nodeList[i].SUSPECTED != true {
						nodeList[i].SUSPECTED = true
						suspectedChan <- nodeList[i]
					}
					if err == nil {
						<-tcpResponseChan
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
				tcp.Send(lead.IP, message.HARTBEATRESPONSE{IP: me.IP})
				fmt.Println("SENT RESPONSE TO", lead.IP)
				time.Sleep(1 * time.Millisecond)
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
			//suspectedChan <- lead
			lead = newL
			nodeList = <-endList
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
