package tcp

import (
	"net"
)

var (
	freeConnections []net.Conn = make([]net.Conn, 0)
)

/*
Method to check if there exists a connection to the address.
If it exists we return this connection. If there does not exist a connection
to the address we construct a new connection to the address and return this.
*/
func Dial(url string) {

}

/*
Method to "close" the connection. We append the connection to the freeConnections
list and return nil.
*/
func Close(conn net.Conn) {

}
