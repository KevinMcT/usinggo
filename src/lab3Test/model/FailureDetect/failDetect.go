package FailureDetect

import (
	"fmt"
	"lab3Test/model/Network/message"
	"lab3Test/model/Network/tcp"
	"time"
)

var (
	delay    = time.NewTicker(2 * time.Second)
	nodelist []message.Node
	errChan  = make(chan error)
	inChan   = make(chan interface{})
	selfnode message.Node
	timer    = time.NewTicker(1 * time.Second)
	timer2   = time.NewTicker(2 * time.Second)
	leader   message.Node
)

func init() {
}

func Fd(nodes chan message.Node, selfnode message.Node) {
	go fillNodelist(nodes, selfnode)
	for {
		if len(nodelist) > 3 {
			if selfnode.LEAD != true {
				lnode, err := recieveHartbeat()
				fmt.Println("Recieved Ping from: ", lnode.IP)
				sendHartbeat(lnode.IP, selfnode)
				fmt.Println("Sent response to: ", lnode.IP)
				if err != nil {
				}
			} else {
				for i, v := range nodelist {
					if v.IP != selfnode.IP {
						fmt.Println("Sending to IP: ", v.IP)
						str := sendHartbeat(v.IP, selfnode)
						slv, err := recieveHartbeat()
						if err != nil {
							fmt.Println(err)
						}
						fmt.Println("Response from: ", slv.IP)
						if str == "suspect" {
							nodelist[i].SUSPECTED = true
							fmt.Println("Suspecting IP: ", v.IP)
						}
						if v.SUSPECTED == true && str == "ok" {
							nodelist[i].SUSPECTED = false
							fmt.Println("IP is back!: ", v.IP)
						}
						<-timer2.C
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
