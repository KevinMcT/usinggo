package machine

import (
	"strconv"
	"strings"
)

var (
	machine     T_Machine
	inputString string
)

type T_Machine struct {
	IP   string
	LEAD bool
	TIME int64
}

func Machine(inputChan chan string, outputChan chan T_Machine) {
	for {
		inputString := <-inputChan
		result := strings.Split(inputString, ":")
		machine.IP = result[0]
		t, _ := strconv.ParseInt(result[1], 10, 64)
		machine.TIME = t
		outputChan <- machine
	}
}
