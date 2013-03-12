package udp

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var ()

func Listen(udpListenChan chan string, udpMasterChan chan string) {
	for {
		mcaddr, _ := net.ResolveUDPAddr("udp4", "239.255.43.99:1889")
		conn, _ := net.ListenMulticastUDP("udp4", nil, mcaddr)
		data := make([]byte, 8192)
		n, addr, _ := conn.ReadFromUDP(data)
		recived := string(data[0:n])
		recivedSplit := strings.Split(recived, ":")
		if strings.Contains(recivedSplit[0], "[B]") {
			SendInitReply("[R]")
			udpListenChan <- addr.IP.String() + ":" + string(recivedSplit[1])
		} else if strings.Contains(recivedSplit[0], "[R]") {
			fmt.Println("Response recived OK from", addr)
			conn.Close()
			break
		}
		time.Sleep(2 * time.Millisecond)
	}
}

func SendBroadcast(startTime int64) {
	mcaddr, _ := net.ResolveUDPAddr("udp4", "239.255.43.99:1889")
	conn, _ := net.ListenMulticastUDP("udp4", nil, mcaddr)
	conn.WriteTo([]byte("[B]:"+strconv.FormatInt(startTime, 10)), mcaddr)
	conn.Close()
}

func SendInitReply(response string) {
	mcaddr, _ := net.ResolveUDPAddr("udp4", "239.255.43.99:1889")
	conn, _ := net.ListenMulticastUDP("udp4", nil, mcaddr)
	conn.WriteTo([]byte(response), mcaddr)
	conn.Close()
}
