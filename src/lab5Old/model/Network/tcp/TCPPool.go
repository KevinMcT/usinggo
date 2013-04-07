package tcp

import (
	"fmt"
	"lab5/Utils"
	"net"
)

type pool struct {
	nConns          int //number of created connections
	freeConnections []net.Conn
}

var instantiated *pool = nil

func init() {
	fmt.Println("Initating TCP pool!")
	instantiated = new(pool)
	instantiated.nConns = 0
	instantiated.freeConnections = make([]net.Conn, 0)
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
			fmt.Println("Error creating connection in TCPPool!!!!!")
			fmt.Println(err)
		}
		instantiated.nConns = instantiated.nConns + 1
		//fmt.Println("Total created connections: ", instantiated.nConns)
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
