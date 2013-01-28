package keyvalueServer

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
)

type Pair struct {
	Key, Value string
}

type Found struct {
	Value string
	Ok    bool
}

type KeyValue map[string]string

func (kv KeyValue) InsertNew(input *Pair, reply *bool) error {
	if _, ok := kv[input.Key]; !ok {
		kv[input.Key] = input.Value
		*reply = true
	} else {
		*reply = false
	}
	return nil
}

func (kv KeyValue) LookUp(input *string, reply *Found) error {
	if v, ok := kv[*input]; ok {
		reply.Value = v
		reply.Ok = true
	} else {
		reply.Ok = false
		reply.Value = ""
	}
	return nil
}

func KeyServer() {
	kv := make(KeyValue)
	rpc.Register(kv)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":12110")
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		rpc.ServeConn(conn)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
