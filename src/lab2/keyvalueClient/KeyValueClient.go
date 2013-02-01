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

func KeyClient(host string) (*rpc.Client, error) {
	service := host

	client, err := rpc.Dial("tcp", service)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	return client, nil
}

func Insert(client *rpc.Client, key string, value string) string {
	args := Pair{key, value}
	var reply bool
	err := client.Call("KeyValue.InsertNew", args, &reply)

	if err != nil {
		log.Fatal("Insert Error:", err)
	}
	return fmt.Sprintf("Inserted: %s with value %s. Result: %v\n", args.Key, args.Value, reply)
}

func LookUp(client *rpc.Client, key string) string {
	var lookup Found
	input := key
	err := client.Call("KeyValue.LookUp", input, &lookup)

	if err != nil {
		log.Fatal("Insert Error:", err)
	}
	var res string
	if lookup.Ok == true {
		res = fmt.Sprintf("Response: LookUp fround value %s at key: %s", lookup.Value, key)
	} else {
		res = fmt.Sprintf("Response: LookUp did not find value %s at key: %s", lookup.Value, key)
	}
	return res
}
