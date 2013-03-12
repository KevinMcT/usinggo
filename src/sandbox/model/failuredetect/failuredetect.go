package failuredetect

import (
	//"fmt"
	"sandbox/controller/machine"
	"sandbox/controller/node"
	"sandbox/model/network/message"
	"sandbox/model/network/tcp"
	"time"
)

var (
	delay int
)

func Detect(me machine.T_Machine, lead node.T_Node, newNodeChan chan node.T_Node, suspectedChan chan node.T_Node, restoreChan chan node.T_Node, tcpRequestChan chan message.HARTBEATREQUEST, tcpResponseChan chan message.HARTBEATRESPONSE, startList []node.T_Node) {
	var ticker = time.NewTicker(500 * time.Millisecond)
	var nodeList = startList
	var newNode node.T_Node
	for {
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(10 * time.Millisecond)
			timeout <- true
		}()
		select {
		case newNode = <-newNodeChan:
			nodeList = AppendIfMissing(nodeList, newNode)
		case <-timeout:
		}
		time.Sleep(1 * time.Millisecond)
		if me.IP == lead.IP {
			for i, v := range nodeList {
				if v.IP != me.IP {
					err := tcp.Send(v.IP, message.HARTBEATREQUEST{IP: lead.IP})
					if err != nil && nodeList[i].SUSPECTED != true {
						nodeList[i].SUSPECTED = true
						suspectedChan <- nodeList[i]
					}
					if err == nil {
						if nodeList[i].SUSPECTED == true {
							nodeList[i].SUSPECTED = false
							restoreChan <- nodeList[i]
						}
						<-tcpResponseChan
					}
					time.Sleep(1 * time.Millisecond)
				}
			}
		}
		if me.IP != lead.IP {
			<-tcpRequestChan
			tcp.Send(lead.IP, message.HARTBEATRESPONSE{IP: me.IP})
			time.Sleep(1 * time.Millisecond)
		}
		<-ticker.C
	}
}

func AppendIfMissing(slice []node.T_Node, i node.T_Node) []node.T_Node {
	for _, ele := range slice {
		if ele.IP == i.IP {
			return slice
		}
	}
	return append(slice, i)
}
