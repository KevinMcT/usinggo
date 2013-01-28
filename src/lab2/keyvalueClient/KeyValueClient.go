package keyvalueClient

import (
	"fmt"
	"log"
	"net/rpc"
)

type Pair struct {
	Key, Value string
}

type Found struct {
	Value string
	Ok    bool
}

func KeyClient(host string) {
	service := host + ":12110"

	client, err := rpc.Dial("tcp", service)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	args := Pair{"1", "Hello"}
	var reply bool
	err = client.Call("KeyValue.InsertNew", args, &reply)

	if err != nil {
		log.Fatal("Insert Error:", err)
	}
	fmt.Printf("Inserted: %s with value %s. Result: %v\n", args.Key, args.Value, reply)

	var lookup Found
	input := "2"
	err = client.Call("KeyValue.LookUp", input, &lookup)

	if err != nil {
		log.Fatal("Insert Error:", err)
	}
	fmt.Printf("Response: %s is %v\n", lookup.Value, lookup.Ok)

}
