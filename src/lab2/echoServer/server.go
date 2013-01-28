// SimpleEchoServer
package main

import (
	"net"
	"os"
	"fmt"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("up4", ":1201")
	checkError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	for {
		if err != nil {
			continue
		}
		handleClient(conn)
		conn.Close() // we're finished
	}
}

func handleClient(conn *net.UDPConn) {
	var buf [512]byte
	for {
		n, addr, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			return
		}
		fmt.Println(string(buf[0:]))
		_, err2 := conn.WriteToUDP(buf[0:n], addr)
		if err2 != nil {
			return
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

