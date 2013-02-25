package main

import (
	"fmt"
	"lab4/model/Client"
)

func main() {
	var tot bool
	tot = true
	for tot == true {
		fmt.Println("---------------------------------------------------------------------")
		fmt.Println("| Choose the function you would like to test using numbers 1 - 4    |")
		fmt.Println("| 1 - Paxos client                                                  |")
		fmt.Println("| 2 - Paxos replica                                                 |")
		fmt.Println("| 0 - Quit                                                          |")
		fmt.Println("---------------------------------------------------------------------")
		var in int
		fmt.Scanf("%d", &in)
		switch in {
		case 1:
			Client.PaxosClient()
		case 2:
			fmt.Println("Not Implemented yet!!!")
		case 0:
			tot = false
		}
	}
}
