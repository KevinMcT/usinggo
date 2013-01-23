package main

import (
	"os"
	"testing/client"
	"testing/server"
)

func main() {
	if os.Args[1] == "server" && os.Args[2] == "tcp" {
		server.TcpServer(os.Args[3])
	}
	if os.Args[1] == "server" && os.Args[2] == "json" {
		server.JsonServer(os.Args[3])
	}

	if os.Args[1] == "client" && os.Args[2] == "tcp" {
		client.TcpClient(os.Args[3])
	}

	if os.Args[1] == "client" && os.Args[2] == "json" {
		client.JsonClient(os.Args[3])
	}

}
