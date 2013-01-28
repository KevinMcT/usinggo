package msgsServer

import (
	"encoding/gob"
	"fmt"
	"lab2/helpers"
	"net"
	"os"
	"os/exec"
)

func MsgsServer(port string) {

	service := "0.0.0.0:" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		decoder := gob.NewDecoder(conn)

		var message helpers.Message
		decoder.Decode(&message)

		conn.Close() // we're finished
	}
}

func demarshal(inChan chan interface{}) {

}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
