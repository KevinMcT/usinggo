package Utils

import (
	"fmt"
	"net"
	"os"
	"reflect"
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

/*
Checks if to structs are the same. Only supports structs that have int and string fields but
not a problem to expand so it can support more types if needed.
*/
func Equals(i1 interface{}, i2 interface{}) bool {
	i1Values := reflect.ValueOf(i1)
	i2Values := reflect.ValueOf(i2)
	fmt.Println("i1 type: ", i1Values.Type())
	fmt.Println("i2 type: ", i2Values.Type())
	if i1Values.Type() == i2Values.Type() {
		fields := i1Values.NumField()
		for i := 0; i < fields; i++ {
			i1Field := i1Values.Field(i)
			i2Field := i2Values.Field(i)
			switch i1Val := i1Field.Interface().(type) {
			case int:
				i2Val := i2Field.Interface().(int)
				if i1Val != i2Val {
					fmt.Println("Not the same int!")
					return false
				}
			case string:
				i2Val := i2Field.Interface().(string)
				if i1Val != i2Val {
					fmt.Println("Not the same string!")
					return false
				}
			}
		}
		return true
	}
	fmt.Println("Not the same type!")
	return false
}
