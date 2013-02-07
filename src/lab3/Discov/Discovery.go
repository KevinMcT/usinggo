package Discov

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
)

func init() {

}

func Listener() (bool, *net.UDPAddr) {
	udpAddr, _ := net.ResolveUDPAddr("udp4", ":2100")
	data := make([]byte, 4096)
	listener, _ := net.ListenUDP("udp4", udpAddr)
	listener.SetDeadline(time.Now().Add(10 * time.Second))
	go Broadcast()
	n, bossIP, err := listener.ReadFromUDP(data)
	if err != nil {
		listener.Close()
		return false, udpAddr
	}
	fmt.Println(string(data[0:n]))
	listener.Close()
	return true, bossIP
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
	fmt.Print("IP: ")
	fmt.Println(BROADCAST_IPv4)
	socket, _ := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   BROADCAST_IPv4,
		Port: port,
	})
	socket.Write([]byte("Found you"))
	socket.Close()
}

func ListenForBroadcast() {
	socket, _ := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 3000,
	})
	socket.SetDeadline(time.Now().Add(30 * time.Second))
	for {
		fmt.Println("Listen!")
		data := make([]byte, 4096)
		_, remoteAddr, err := socket.ReadFromUDP(data)
		if err != nil {
			break
		}
		socket.SetDeadline(time.Now().Add(30 * time.Second))
		ips = append(ips, remoteAddr)
		fmt.Println(ips)
		fmt.Println("Responding....")
		RespondToBroadcast()
		fmt.Println("Response sent....")
	}
}
