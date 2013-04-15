package node

import (
	"strconv"
	"strings"
)

var (
	node T_Node
)

type T_Node struct {
	IP        string
	ROLE      string
	TIME      int64
	SUSPECTED bool
	ALIVE     bool
	LEAD      bool
}

func Node(inputChan chan string, outputChan chan T_Node) {
	for {
		inputString := <-inputChan
		result := strings.Split(inputString, ":")
		node.IP = result[0]
		t, _ := strconv.ParseInt(result[1], 10, 64)
		node.TIME = t
		node.ALIVE = true
		outputChan <- node
	}
}
