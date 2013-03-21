package FailureDetect

import (
	"fmt"
	"lab5/model/Network/message"
	"lab5/model/Network/tcp"
	"lab5/model/RoundVar"
	"time"
)

var (
	nodelist []message.Node
	errChan  = make(chan error)
	inChan   = make(chan interface{})
	selfnode message.Node
	timer    = time.NewTicker(300 * time.Millisecond)
	leader   message.Node
)

func init() {
}

func Fd(nodes chan message.Node, selfnode message.Node, leadElect chan []message.Node) {
	go fillNodelist(nodes, selfnode)
	var suspectedNode message.Node
	for {
		if len(nodelist) > 1 {
			if selfnode.LEAD != true {
				var myLead message.Node
				for i, v := range nodelist {
					if v.LEAD == true {
						myLead = nodelist[i]
					}
				}
				if myLead.IP != "" {
					lnode, err := recieveHartbeat()
					if lnode.SUSPECTED == true && lnode.IP == "" && err != nil {
						for i, v := range nodelist {
							if v.IP == myLead.IP {
								nodelist[i].SUSPECTED = true
							}
							nodelist[i].LEAD = false
						}
						leadElect <- nodelist
						nodelist = make([]message.Node, 0)
						myLead = message.Node{}
						break
					}
					if err == nil {
						if lnode.SUSPECTED == true {
							for i, v := range nodelist {
								if lnode.IP == v.IP {
									nodelist[i].SUSPECTED = true
								}
							}
						} else {
							for i, v := range nodelist {
								if lnode.IP == v.IP {
									if v.SUSPECTED == true {
										nodelist[i].SUSPECTED = false
									} else {

									}
								}
							}
						}
						sendHartbeat(myLead.IP, selfnode)
					}
				}
			} else {
				for i, v := range nodelist {
					if v.IP != selfnode.IP {
						str := sendHartbeat(v.IP, suspectedNode)
						if str == "ok" {
							_, err := recieveHartbeat()
							if err != nil {
								fmt.Println(err)
							}
						}
						if str == "suspect" {
							nodelist[i].SUSPECTED = true
							suspectedNode = nodelist[i]
						}
						if v.SUSPECTED == true && str == "ok" {
							nodelist[i].SUSPECTED = false
							suspectedNode = nodelist[i]
						}
					}
				}
			}
		}
		<-timer.C
	}
}

func sendHartbeat(ip string, node message.Node) string {
	err := tcp.Send(ip, node)
	if err != nil {
		return "suspect"
	}
	if node.SUSPECTED == true && err == nil {
		return "ok"
	}
	return "ok"
}

func recieveHartbeat() (message.Node, error) {
	node, err := tcp.Recieve()
	return node, err
}

func fillNodelist(nc chan message.Node, self message.Node) {
	if self.LEAD == false {
		nodelist = append(nodelist, self)
	}
	if self.LEAD == true {
		leader = self
	}
	for {
		node := <-nc
		if node.LEAD == true {
			leader = node
		}
		fmt.Println("added: ", node)
		nodelist = append(nodelist, node)
		RoundVar.GetRound().List = nodelist
	}
}
