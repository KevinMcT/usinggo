package Utils

import (
	"fmt"
	"net"
	"os"
	"strings"
)

/*
	Utils class where we try to collect all utils that are used for multiple classes
*/

/*
	Checks systems errors that will corrupt the whole system.
	err 	error 	error to be checked and printed
*/
func CheckError(err error) {
	if err != nil {
		fmt.Println("--Fatal error", err.Error(), "--")
		os.Exit(1)
	}
}

/*
	Searches for a specified IP in an array and returns a connection for that IP
	Is beeing used in the TCPPool class
	IP		string		IP to be searched for
	array	[]net.Conn	The array to search through
	returns net.Conn	Connection
*/
func SearchForIP(IP string, array []net.Conn) net.Conn {
	for _, v := range array {
		if v.RemoteAddr().String() == IP {
			return v.(net.Conn)
		}
	}
	return nil
}

/*
	Whenever we need the IP without the port.
	address 	string		The connection
	returns		string		IP without port
*/
func GetIp(address string) string {
	remoteSplit := strings.Split(address, ":")
	return remoteSplit[0]
}
