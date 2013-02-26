package main

import (
	"fmt"
	"sandbox/controller/machine"
	"sandbox/model/network/udp"
	"time"
)

var (
	udpListenChan = make(chan string, 10)
)

func main() {
	fmt.Println("Go")
	startTime := time.Now().UnixNano()
	fmt.Println("Start Listen")
	go udp.Listen(udpListenChan)
	fmt.Println("Send broadcast")
	udp.SendBroadcast(startTime)
	go machine.Machine(udpListenChan)
	for {
		time.Sleep(100 * time.Millisecond)
	}
}
