package main

import (
	"fmt"
	"lab2/echoClient"
	"lab2/echoServer"
	"lab2/keyvalueClient"
	"lab2/keyvalueServer"
	"lab2/msgsClient"
	"lab2/msgsServer"
)

func main() {
	fmt.Println("---------------------------------------------------------------------")
	fmt.Println("| Choose the function you would like to test using numbers 1 - 4    |")
	fmt.Println("| 1 - TCP Server                                                    |")
	fmt.Println("| 2 - TCP Client                                                    |")
	fmt.Println("| 3 - RPC Server                                                    |")
	fmt.Println("| 4 - RPC Client                                                    |")
	fmt.Println("| 5 - Echo Server                                                   |")
	fmt.Println("| 6 - Echo Client                                                   |")
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
	case 3:
		fmt.Println("Starting rcp server")
		keyvalueServer.KeyServer()
	case 4:
		fmt.Println("Enter hostname")
		var rpc string
		fmt.Scanf("%s", &rpc)
		keyvalueClient.KeyClient(rpc)
	case 5:
		fmt.Println("Starting server")
		echoServer.EchoServer()
	case 6:
		fmt.Println("Enter hostname")
		var echo string
		fmt.Scanf("%s", &echo)
		echoClient.EchoClient(echo)
	}
}
