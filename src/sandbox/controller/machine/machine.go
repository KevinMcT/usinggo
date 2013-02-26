package machine

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	machine     t_Machine
	inputString string
)

type t_Machine struct {
	IP   string
	ROLE string
	TIME int64
}

func Machine(inputChan chan string) {
	for {
		inputString := <-inputChan
		result := strings.Split(inputString, ":")
		machine.IP = result[0]
		t, _ := strconv.ParseInt(result[1], 10, 64)
		machine.TIME = t
		fmt.Println(machine)
	}
}
