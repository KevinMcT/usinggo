package tcp

import (
	"encoding/gob"
	"fmt"
	"lab3Test/model/Network/message"
	"net"
	"os"
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
	service := ip + ":2000"
	conn, err := net.Dial("tcp", service)
	fmt.Println(err)
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
	service := "0.0.0.0:2000"
	tcpAddr, _ := net.ResolveTCPAddr("tcp", service)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	conn, err2 := listener.Accept()
	var msg interface{}
	if err != nil {
		fmt.Println("Recieve. 43: ", err)
	}
	if err2 != nil {
		fmt.Println("Recieve. 46: ", err2)
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

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error", err.Error())
		os.Exit(1)
	}
}
