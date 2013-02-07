package main

import (
	"fmt"
	"lab3/Discov"
)

var (
	end = make(chan int)
)

func main() {
	ok, ip := Discov.Listener()
	if ok {
		fmt.Println("Other process is master")
		fmt.Print("Boss IP: ")
		fmt.Println(ip)
	} else {
		fmt.Println("I am master")
		Discov.ListenForBroadcast()
		<-end
	}
}
