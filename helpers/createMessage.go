package helpers

import (
	"encoding/binary"
)

func CreateMessage( op uint, msg string) []byte {

	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, uint32(len(msg)+1))
	// fmt.Println(sizeBytes)

	var buffer []byte
	buffer = append(buffer, append(sizeBytes, append([]byte{byte(op)}, []byte(msg)...)...)...)

	// fmt.Println(buffer)
	return buffer

}
