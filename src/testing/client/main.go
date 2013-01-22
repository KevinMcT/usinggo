package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"testing/helpers"
)

func main() {

	person := helpers.Person{Name: "Patrik", Adress: "Solåsveien 34", Mail: "psb@psb.no"}

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port", os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1]

	conn, err := net.Dial("tcp", service)
	checkError(err)

	encoder := gob.NewEncoder(conn)
	encoder.Encode(person)

	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
