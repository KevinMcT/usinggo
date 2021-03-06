package main

import (
	"fmt"
	"lab4/model/Client"
	"lab4/model/Replica"
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
			fmt.Println("Starting Paxos client...")
			Client.PaxosClient()
			break
		case 2:
			fmt.Println("Starting Paxos replica...")
			Replica.PaxosReplica()
			break
		case 0:
			tot = false
			break
		}
	}
}
