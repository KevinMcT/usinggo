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
	for _, v := range array {
		if v.RemoteAddr().String() == IP {
			return v.(net.Conn)
		}
	}
	return nil
}
