package main

import (
	"fmt"
	"lab6/core"
	"os"
)

var ()

func main() {
	var tot bool
	tot = true
	for tot == true {
		fmt.Println("*********************************************************************")
		fmt.Println("* Choose the function you would like to test using numbers 1 - 2    *")
		fmt.Println("* 1 - Paxos server                                                  *")
		fmt.Println("* 2 - Client                                                        *")
		fmt.Println("* 3 - Bank system													 *")
		fmt.Println("* 0 - Quit                                                          *")
		fmt.Println("*********************************************************************")
		var in int
		fmt.Scanf("%d", &in)

		switch in {
		case 1:
			//--Calling the Paxos replica--//
			core.Paxos()
		case 2:
			//--Calling the Client--//
			core.Client()
		case 3:
			core.Bank()
		case 0:
			os.Exit(0)
		}
	}
}
