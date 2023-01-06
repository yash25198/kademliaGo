package network

import (
	"fmt"
	"net"
)

func Send(conn net.Conn, msg []byte) bool {
	_, errSize := conn.Write(msg)
	if errSize != nil {
		fmt.Println("Error writing :", errSize.Error())
		return false
	}
	return true
}
