package helpers

import (
	"fmt"
	"net"
)

func CreateConnection(address string) net.Conn {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// fmt.Println(conn.LocalAddr().String())
	return conn
}
