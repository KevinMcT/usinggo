package Paxos

import (
	"encoding/gob"
	"fmt"
	//lab4/Utils"
	"lab4/model/Network/message"
	"net"
)

var (
	leader   message.Node
	self     message.Node
	nodeList = make([]message.Node, 0)
	round    int
)

func Acceptor(led message.Node, me message.Node, nc chan message.Node, ac chan string) {
	round = 0
	self = me
	go fillNodelist(nc)
	go receviedPromise()
	for {
		value := <-ac
		fmt.Println("Received value ", value)
		round = round + 1
		sendValue(value)
	}
}

func fillNodelist(nc chan message.Node) {
	if self.LEAD == false {
		nodeList = append(nodeList, self)
	}
	if self.LEAD == true {
		leader = self
	}
	for {
		node := <-nc
		if node.LEAD == true {
			leader = node
		}
		nodeList = append(nodeList, node)
	}
}

func sendValue(value string) {
	fmt.Println(nodeList)
	for _, v := range nodeList {

		//Send prepare
		sendAddress := v.IP + ":1338"
		fmt.Println("Sending prepare to ", sendAddress)
		sendConn, err := net.Dial("tcp", sendAddress)
		if err != nil {
			fmt.Println(err)
		} else {
			encoder := gob.NewEncoder(sendConn)
			var prepare = message.Prepare{ROUND: round}
			var msg interface{}
			msg = prepare
			encoder.Encode(&msg)
		}
		sendConn.Close()
		/*fmt.Println("Sendt msg, waiting for promise")
		//Wait for promise 
		receiveAddress := "0.0.0.0:8761"
		tcpAddr, err := net.ResolveTCPAddr("tcp", receiveAddress)
		Utils.CheckError(err)
		listener, err := net.ListenTCP("tcp", tcpAddr)
		Utils.CheckError(err)
		receiveConn, err := listener.Accept()
		decoder := gob.NewDecoder(receiveConn)
		var msg interface{}
		err = decoder.Decode(&msg)
		if err != nil {
			fmt.Println(err)
		}
		if msg != nil {
			var promise message.Promise
			promise = msg.(message.Promise)
			fmt.Println("Received promise ", promise)
		}

		//Send accept
		fmt.Println("Sending accept")
		encoder := gob.NewEncoder(sendConn)
		var accept = message.Accept{ROUND: round, VALUE: value}
		msg = accept
		encoder.Encode(&msg)

		fmt.Println("Closing channels...")
		receiveConn.Close()
		fmt.Println("Channels closed...")*/
	}
}

func receviedPromise() {
	for {
		value := <-message.PromiseChan
		address := value.Ip + ":1338"
		fmt.Println("Sending accept to ", address)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Println(err)
		} else {
			encoder := gob.NewEncoder(conn)
			var accept = message.Accept{ROUND: round, VALUE: "test"}
			var msg interface{}
			msg = accept
			encoder.Encode(&msg)
		}
		conn.Close()
	}
}
