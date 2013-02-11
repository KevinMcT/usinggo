package LeaderElect

import (
	"fmt"
	"lab3/messages"
	"net"
	"strconv"
)

var (
	nodes = make([]messages.Node, 0)
)

func Elect(mytime string, _nodes []messages.Node) *net.UDPAddr {
	mysTime, _ := strconv.Atoi(mytime)
	hissTime, _ := strconv.Atoi(_nodes[1].Time)
	if mysTime < hissTime {
		fmt.Print("Me older!")
		fmt.Println(mysTime)
		fmt.Println(hissTime)
		return nil
	}
	fmt.Println("He older!")
	fmt.Println(mysTime)
	fmt.Println(hissTime)
	return _nodes[1].IP
}
