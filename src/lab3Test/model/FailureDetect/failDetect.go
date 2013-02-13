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
				node, err := recieveHartbeat()
				fmt.Println(node)
				fmt.Println("Moving on")
				if err != nil {
					fmt.Println("WTF!")
				}
			} else {
				for _, v := range nodelist {
					if v.IP != selfnode.IP {
						fmt.Println("Sending to IP: ", v.IP)
						sendHartbeat(v.IP, selfnode)
						<-timer2.C
					}
				}
			}
		}
		<-timer.C
	}
}

func sendHartbeat(ip string, node message.Node) {
	var err2 error
	fmt.Println("Hartbeat")
	err := tcp.Send(ip, node)
	if err != nil {
		err2 = tcp.Send(ip, node)
		for err2 != nil {
			err2 = tcp.Send(ip, node)
		}
	}
}

func recieveHartbeat() (message.Node, error) {
	fmt.Println("Recieve")
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
