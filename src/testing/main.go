package main

import (
	"fmt"
	"testing/client"
	"testing/server"
)

func main() {
	fmt.Println("---------------------------------------------------------------------")
	fmt.Println("| Choose the function you would like to test using numbers 1 - 4    |")
	fmt.Println("| 1 - TCP Server                                                    |")
	fmt.Println("| 2 - json Server                                                   |")
	fmt.Println("| 3 - TCP Client                                                    |")
	fmt.Println("| 4 - json Client                                                   |")
	fmt.Println("| 0 - Quit                                                          |")
	fmt.Println("---------------------------------------------------------------------")
	var in int
	fmt.Scanf("%d", &in)

	switch in {
	case 1:
		fmt.Println("Enter port number: ")
		var port string
		fmt.Scanf("%s", &port)
		server.TcpServer(port)
	case 2:
		fmt.Println("Enter port number: ")
		var port string
		fmt.Scanf("%s", &port)
		server.JsonServer(port)
	case 3:
		fmt.Println("Enter host and port eg 'localhost:port'")
		var host string
		fmt.Scanf("%s", &host)
		client.TcpClient(host)
	case 4:
		fmt.Println("Enter host and port eg 'localhost:port'")
		var host string
		fmt.Scanf("%s", &host)
		client.JsonClient(host)
	}
}
