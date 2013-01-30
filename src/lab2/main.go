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
	var tot bool
	tot = true
	for tot == true {
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
			fmt.Println("Starting TCP server")
			msgsServer.MsgsServer(port)
		case 2:
			fmt.Println("Enter host and port eg 'localhost:port': ")
			var host string
			fmt.Scanf("%s", &host)
			fmt.Println("Starting TCP client")
			msgsClient.MsgsClient(host)
		case 3:
			fmt.Println("Enter port number: ")
			var rpcPort string
			fmt.Scanf("%s", &rpcPort)
			fmt.Println("Starting RCP server")
			keyvalueServer.KeyServer(rpcPort)
		case 4:
			fmt.Println("Enter hostname and port eg 'localhost:port': ")
			var rpc string
			fmt.Scanf("%s", &rpc)
			fmt.Println("Starting RCP client")
			client, _ := keyvalueClient.KeyClient(rpc)
			var run bool
			run = true
			for run == true {
				fmt.Println("---------------------------------------------------------------------")
				fmt.Println("| Insert or LookUp                                                  |")
				fmt.Println("| 1 - Insert                                                        |")
				fmt.Println("| 2 - LookUp                                                        |")
				fmt.Println("| 0 - Exit                                                          |")
				fmt.Println("---------------------------------------------------------------------")
				var rpcChoice int
				fmt.Scanf("%d", &rpcChoice)
				switch rpcChoice {
				case 1:
					fmt.Println("Enter key to insert: ")
					var key string
					fmt.Scanf("%s", &key)
					fmt.Println("Enter value to insert on that key: ")
					var value string
					fmt.Scanf("%s", &value)
					res := keyvalueClient.Insert(client, key, value)
					fmt.Println(res)
				case 2:
					var key string
					fmt.Println("Enter key to lookup: ")
					fmt.Scanf("%s", &key)
					res := keyvalueClient.LookUp(client, key)
					fmt.Println(res)
				case 0:
					run = false
				}
			}
		case 5:
			fmt.Println("Enter port number: ")
			var echoPort string
			fmt.Scanf("%s", &echoPort)
			fmt.Println("Starting echo server")
			echoServer.EchoServer(echoPort)
		case 6:
			fmt.Println("Enter hostname and port eg 'localhost:port': ")
			var echo string
			fmt.Scanf("%s", &echo)
			echoClient.EchoClient(echo)
		case 0:
			tot = false
		}
	}
}
