package main

import (
	"btb/core"
	"fmt"
	"os"
)

var ()

func main() {
	var tot bool
	tot = true
	for tot == true {
		fmt.Println("---------------------------------------------------------------------")
		fmt.Println("| Choose the function you would like to test using numbers 1 - 2    |")
		fmt.Println("| 1 - Paxos server                                                  |")
		fmt.Println("| 2 - Client                                                        |")
		fmt.Println("| 0 - Quit                                                          |")
		fmt.Println("---------------------------------------------------------------------")
		var in int
		fmt.Scanf("%d", &in)

		switch in {
		case 1:
			core.Paxos()
		case 2:
			fmt.Println("--Call client here--")
		case 0:
			os.Exit(0)
		}
	}
}
