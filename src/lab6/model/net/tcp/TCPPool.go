package tcp

import (
	"encoding/gob"
	"fmt"
	"lab6/Utils"
	"net"
)

type pool struct {
	nConns          int //number of created connections
	freeConnections []net.Conn
	test            map[net.Conn]*gob.Encoder
}

var instantiated *pool = nil

func init() {
	fmt.Println("--Initating TCP pool--")
	instantiated = new(pool)
	instantiated.nConns = 0
	instantiated.freeConnections = make([]net.Conn, 0)
	instantiated.test = make(map[net.Conn]*gob.Encoder)
}

func GetEncoder(url string) *gob.Encoder {
	for i, _ := range instantiated.test {
		if i.RemoteAddr().String() == url {
			encoder := instantiated.test[i]
			return encoder
		}
	}
	conn, _ := net.Dial("tcp", url)
	encoder := gob.NewEncoder(conn)
	return encoder
}

func StoreEncoder(conn net.Conn, encoder gob.Encoder) *gob.Encoder {
	for i, _ := range instantiated.test {
		if i == conn {
			return nil
		}
	}
	instantiated.test[conn] = &encoder
	return nil
}

/*
Method to check if there exists a connection to the address.
If it exists we return this connection. If there does not exist a connection
to the address we construct a new connection to the address and return this.
*/
func Dial(url string) net.Conn {
	conn := Utils.SearchForIP(url, instantiated.freeConnections)
	if conn == nil {
		conn, err := net.Dial("tcp", url)
		if err != nil {
			fmt.Println("--Error creating connection in TCPPool:", err, "--")
		}
		instantiated.nConns = instantiated.nConns + 1
		return conn
	}
	return conn
}

/*
Method to "close" the connection. We append the connection to the freeConnections
list and return nil.
*/
func Close(conn net.Conn) net.Conn {
	existConn := Utils.SearchForIP(conn.RemoteAddr().String(), instantiated.freeConnections)

	//If the connection is allready in the list we don`t add it again.
	if existConn == nil {
		instantiated.freeConnections = append(instantiated.freeConnections, conn)
	}
	return nil
}
