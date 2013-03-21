package Utils

import (
	"fmt"
	"net"
	"os"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error", err.Error())
		os.Exit(1)
	}
}

func SearchForIP(IP string, array []net.Conn) net.Conn {
	fmt.Println("IP:", IP)
	for _, v := range array {
		//fmt.Println("connIP:", v.RemoteAddr().String())
		if v.RemoteAddr().String() == IP {
			//fmt.Println("Found a valid connection in the")
			return v.(net.Conn)
		}
	}
	return nil
}
