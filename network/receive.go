package network
import (
	"fmt"
	"net"
)
func Receive(conn net.Conn,size uint) []byte {
	buf := make([]byte, size)
	_, errSize := conn.Read(buf)
	if errSize != nil {
		fmt.Println("Error reading :", errSize.Error())
	}
	return buf
}
