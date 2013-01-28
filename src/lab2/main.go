package main

import (
	"fmt"
	"lab2/msgsClient"
	"lab2/msgsServer"
)

func main() {
	fmt.Println("---------------------------------------------------------------------")
	fmt.Println("| Choose the function you would like to test using numbers 1 - 4    |")
	fmt.Println("| 1 - Server                                                        |")
	fmt.Println("| 2 - Client                                                        |")
	fmt.Println("| 0 - Quit                                                          |")
	fmt.Println("---------------------------------------------------------------------")
	var in int
	fmt.Scanf("%d", &in)

	switch in {
	case 1:
		fmt.Println("Enter port number: ")
		var port string
		fmt.Scanf("%s", &port)
		msgsServer.MsgsServer(port)
	case 2:
		fmt.Println("Enter host and port eg 'localhost:port'")
		var host string
		fmt.Scanf("%s", &host)
		msgsClient.MsgsClient(host)
	}
}
