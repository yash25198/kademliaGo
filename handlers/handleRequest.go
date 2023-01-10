// Handles incoming requests.
package handlers

import (
	"encoding/binary"
	"fmt"
	"kademliago/helpers"
	"kademliago/kademlia"
	"kademliago/network"
	"kademliago/types"
	"net"
	"strings"
)

func HandleRequest(conn net.Conn, rt *kademlia.RoutingTable, multiaddr string, msgqueue chan types.Msgq) {
	buf := network.Receive(conn, 4)
	size := binary.BigEndian.Uint32(buf)

	// fmt.Println(size)

	buf2 := network.Receive(conn, uint(size))

	opcode := buf2[0]
	switch opcode {
	case 1:
		{
			data := string(buf2[1:])
			fmt.Println("ping data :", data)

			split := strings.Split(data, "/")
			// fmt.Println(split)

			ip := split[3]
			port := strings.Split(ip, ":")
			nodeid := split[1]
			//extract ip
			// fmt.Println(ip, port[1], nodeid)
			exconn := helpers.CreateConnection(ip)
			if exconn == nil {
				break
			}
			msg := helpers.CreateMessage(2, multiaddr)

			network.Send(exconn, msg)
			rt.AddNode(kademlia.Node{
				IP:     port[0],
				Port:   port[1],
				NodeID: nodeid,
			})
		}
	case 2:
		{
			data := string(buf2[1:])
			fmt.Println("pong data :", data)

			split := strings.Split(data, "/")

			ip := split[3]
			port := strings.Split(ip, ":")
			nodeid := split[1]

			rt.AddNode(kademlia.Node{
				IP:     port[0],
				Port:   port[1],
				NodeID: nodeid,
			})
		}
	case 3:
		{

			data := string(buf2[1:])
			split2 := strings.Split(data, "+")
			split := strings.Split(split2[0], "/")
			ipport := strings.Split(split[3], ":")
			rt.AddNode(kademlia.Node{
				NodeID: split[1],
				IP:     ipport[0],
				Port:   ipport[1],
			})
			fmt.Println("text data :", split2[1])

		}
	case 4:
		{
			data := string(buf2[1:])
			fmt.Println("find node data :", data)

			split2 := strings.Split(data, "+")
			split := strings.Split(split2[0], "/")

			payload := split2[1]
			ipip := strings.Split(split[3], ":")
			ip := ipip[0]
			port := ipip[1]

			exconn := helpers.CreateConnection(ip + ":" + port)
			if exconn == nil {
				break
			}

			nodeID, ipFound, port, found := rt.FindNode(payload, split[1])
			if found {
				msg := helpers.CreateMessage(6, multiaddr+"+/"+nodeID+"/tcp/"+ipFound+":"+port)
				network.Send(exconn, msg)
			} else {
				msg := helpers.CreateMessage(5, multiaddr+"+/"+nodeID+"/tcp/"+ipFound+":"+port+"+"+payload)
				network.Send(exconn, msg)
			}
			// rt.AddNode(kademlia.Node{
			// 	IP:     ip,
			// 	Port:   port,
			// 	NodeID: split[1],
			// })

		}
	case 5:
		{
			data := string(buf2[1:])
			fmt.Println("5 data :", data)

			split2 := strings.Split(data, "+")
			split := strings.Split(split2[1], "/")

			ipsplit := strings.Split(split[3], ":")
			ip := ipsplit[0]
			port := ipsplit[1]
			node_id := split[1]

			exconn := helpers.CreateConnection(ip + ":" + port)
			if exconn == nil {
				break
			}
			msg := helpers.CreateMessage(4, multiaddr+"+"+split2[2])
			network.Send(exconn, msg)
			rt.AddNode(kademlia.Node{
				IP:     ip,
				Port:   port,
				NodeID: node_id,
			})

		}
	case 6:
		{
			data := string(buf2[1:])
			fmt.Println("6 data :", data)

			split2 := strings.Split(data, "+")
			split := strings.Split(split2[1], "/")

			ipip := strings.Split(split[3], ":")

			rt.AddNode(kademlia.Node{
				NodeID: split[1],
				Port:   ipip[1],
				IP:     ipip[0],
			})

			msgqueue <- types.Msgq{
				Receiver: split[1],
				IP:       split[3],
			}

		}
	}

	conn.Close()

}
