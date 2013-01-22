package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/exec"
	"testing/helpers"
)

func main() {

	service := "0.0.0.0:1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		decoder := gob.NewDecoder(conn)

		var person helpers.Person
		decoder.Decode(&person)
		if person.Name == "Patrik" {
			fmt.Println("Boss is here!")
			fmt.Println(person.Name + ":" + person.Adress + ":" + person.Mail)
			cmd := exec.Command("firefox", "http://rotor.no")
			err := cmd.Start()
			werr := cmd.Wait()
			conn.Write([]byte("Command closed"))
			if err != nil {
				fmt.Println(err)
			}
			if werr != nil {
				fmt.Println(werr)
			}

		} else {
			fmt.Println("You are: " + person.Name)
			cmd := exec.Command("firefox", "http://9gag.com")
			err := cmd.Start()
			werr := cmd.Wait()
			conn.Write([]byte("Command closed"))
			if err != nil {
				fmt.Println(err)
			}
			if werr != nil {
				fmt.Println(werr)
			}
		}

		conn.Close() // we're finished
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
