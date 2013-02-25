package FailureDetect

import (
	"fmt"
	"lab4/model/Network/message"
	"lab4/model/Network/tcp"
	"time"
)

var (
	nodelist []message.Node
	errChan  = make(chan error)
	inChan   = make(chan interface{})
	selfnode message.Node
	timer    = time.NewTicker(5 * time.Second)
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
						fmt.Println("Recieved Ping from: ", myLead.IP)
						if lnode.SUSPECTED == true {
							for i, v := range nodelist {
								if lnode.IP == v.IP {
									nodelist[i].SUSPECTED = true
									fmt.Println("Leader says suspected on ", v.IP)
								}
							}
						} else {
							for i, v := range nodelist {
								if lnode.IP == v.IP {
									if v.SUSPECTED == true {
										nodelist[i].SUSPECTED = false
										fmt.Println("Leader says ok on ", v.IP)
									} else {

									}
								}
							}
						}
						str := sendHartbeat(myLead.IP, selfnode)
						if str == "ok" {
							fmt.Println("Sent response to: ", myLead.IP)
						}
					}
				}
			} else {
				for i, v := range nodelist {
					if v.IP != selfnode.IP {
						fmt.Println("Sending to IP: ", v.IP)
						str := sendHartbeat(v.IP, suspectedNode)
						if str == "ok" {
							slv, err := recieveHartbeat()
							if err != nil {
								fmt.Println(err)
							}
							fmt.Println("Response from: ", slv.IP)
						}
						if str == "suspect" {
							nodelist[i].SUSPECTED = true
							fmt.Println("Suspecting IP: ", v.IP)
							suspectedNode = nodelist[i]
						}
						if v.SUSPECTED == true && str == "ok" {
							nodelist[i].SUSPECTED = false
							fmt.Println("IP is back!: ", v.IP)
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
	}
}
