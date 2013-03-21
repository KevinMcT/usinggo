package tcp

import (
	"encoding/gob"
	"fmt"
	"lab5/model/Network/message"
	"net"
	"time"
)

var (
	tick = time.NewTicker(2 * time.Second)
)

func init() {

}

func Send(ip string, node message.Node) error {
	var msg interface{}
	msg = node
	service := ip + ":2001"
	conn, err := net.Dial("tcp", service)
	if err != nil {
		return err
	} else {
		encoder := gob.NewEncoder(conn)
		encoder.Encode(&msg)
	}
	conn.Close()
	return err
}

func Recieve() (message.Node, error) {
	var node message.Node
	service := "0.0.0.0:2001"
	tcpAddr, _ := net.ResolveTCPAddr("tcp", service)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	listener.SetDeadline(time.Now().Add(1900 * time.Millisecond))
	conn, err2 := listener.Accept()
	var msg interface{}
	if err != nil {
		node = message.Node{SUSPECTED: true}
		return node, err2
	}
	if err2 != nil {
		node = message.Node{SUSPECTED: true}
		listener.Close()
		return node, err2
	}
	decoder := gob.NewDecoder(conn)
	errDec := decoder.Decode(&msg)
	if errDec != nil {
		fmt.Println(errDec)
	}
	if msg != nil {
		node = msg.(message.Node)
	}
	listener.Close()
	conn.Close()
	return node, err
}
