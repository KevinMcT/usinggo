package FD

import (
	"encoding/gob"
	"fmt"
	"lab3/messages"
	"net"
	"time"
)

var ()

func ListenForUDP() {
	fmt.Println("Listen routine!")
	udpAddr, _ := net.ResolveUDPAddr("udp4", ":2500")
	listener, _ := net.ListenUDP("udp4", udpAddr)
	listener.SetDeadline(time.Now().Add(15 * time.Second))
	for {
		data := make([]byte, 4096)
		n, _, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Println("FAILURE! LEADER MIGHT BE DOWN!")
			fmt.Println(err)
			break
		}
		fmt.Println(string(data[0:n]))
		listener.SetDeadline(time.Now().Add(20 * time.Second))
	}
	listener.Close()
}

func WriteUDP(message messages.Ping, addr *net.UDPAddr) {
	hostname := addr.IP.String()
	hostname = hostname + ":2500"
	resvAddr, _ := net.ResolveUDPAddr("udp4", hostname)
	socket, _ := net.DialUDP("udp4", nil, resvAddr)
	encoder := gob.NewEncoder(socket)
	encoder.Encode(message)
	socket.Close()
}
