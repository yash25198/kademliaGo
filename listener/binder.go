package listener

import (
	// "encoding/binary"
	"fmt"
	"kademliago/handlers"
	"kademliago/kademlia"
	"kademliago/types"
	"net"
	"os"
)

func BindServer(CONN_HOST string, CONN_PORT string, CONN_TYPE string, NODE_ID string, rt *kademlia.RoutingTable, msgqueue chan types.Msgq) {
	// Listen for incoming connections

	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	MUL_ADDR := "/" + NODE_ID + "/" + CONN_TYPE + "/" + CONN_HOST + ":" + CONN_PORT
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		handlers.HandleRequest(conn, rt, MUL_ADDR, msgqueue)

	}
}
