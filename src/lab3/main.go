package main

import (
	"fmt"
	"net"
	"time"
)

var (
	port   = 2100
	listen = make(chan int)
	ticker = time.NewTicker(time.Second * 30)
	ips    = make([]*net.UDPAddr, 0)
	rec    = make(chan int)
	end    = make(chan int)
)

func main() {
	ok := Listener()
	if ok {
		fmt.Println("Other process is master")
	} else {
		fmt.Println("I am master")
		ListenForBroadcast()
		<-end
	}
}

func Listener() bool {
	var ok bool
	udpAddr, _ := net.ResolveUDPAddr("up4", ":2100")
	data := make([]byte, 4096)
	listener, _ := net.ListenUDP("udp4", udpAddr)
	go Broadcast()
	listener.SetDeadline(time.Now().Add(5 * time.Second))
	for {
		n, err := listener.Read(data)
		if err != nil {
			ok = false
			break
		}
		fmt.Println(string(data[0:n]))
		ok = true
	}
	return ok
}

func Broadcast() {
	BROADCAST_IPv4 := net.IPv4(255, 255, 255, 255)
	socket, _ := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   BROADCAST_IPv4,
		Port: 3000,
	})
	socket.Write([]byte("broadcast"))
	fmt.Println("Broadcasted myself")
	socket.Close()
}

func RespondToBroadcast() {
	BROADCAST_IPv4 := ips[len(ips)-1].IP
	socket, _ := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   BROADCAST_IPv4,
		Port: 2100,
	})
	socket.Write([]byte("Found you"))
	socket.Close()
}

func ListenForBroadcast() {
	socket, _ := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: port,
	})
	socket.SetDeadline(time.Now().Add(30 * time.Second))
	for {
		fmt.Println("Listen!")
		data := make([]byte, 4096)
		read, remoteAddr, err := socket.ReadFromUDP(data)
		if err != nil {
			break
		}
		socket.SetDeadline(time.Now().Add(30 * time.Second))
		fmt.Println(string(data[0:read]))
		fmt.Println(remoteAddr.IP)
		ips = append(ips, remoteAddr)
		rec <- 1
		RespondToBroadcast()
	}
}
