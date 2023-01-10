package main

import (
	"fmt"
	"kademliago/helpers"
	"kademliago/kademlia"
	"kademliago/listener"
	"kademliago/network"
	"kademliago/types"
	"os"
)

func main() {
	args := os.Args
	rt := kademlia.CreatRoutingTable(160, args[4])
	ch := make(chan types.Msgq)
	go listener.BindServer(args[1], args[2], args[3], args[4], &rt, ch)
	var addr string
	var op uint
	var msg string
	for {

		fmt.Scanf("%s %d %s", &addr, &op, &msg)
		if !(len(addr) == 40) {
			conn := helpers.CreateConnection(addr)
			encodedmsg := helpers.CreateMessage(op, msg)
			stat := network.Send(conn, encodedmsg)
			if stat {
				fmt.Println("msg sent")
			}
		} else {
			if op == 3 {
				_, ip, port, found := rt.FindNode(addr, args[4])
				if !found && ip != "NULL" {
					// fmt.Println("herre",ip,port)
					conn := helpers.CreateConnection(ip + ":" + port)
					encodedmsg := helpers.CreateMessage(4, "/"+args[4]+"/"+args[3]+"/"+args[1]+":"+args[2]+"+"+addr)
					stat := network.Send(conn, encodedmsg)
					if stat {
						fmt.Println("request sent")
					}
					m := types.Msgq{}
					m = <-ch
					conn2 := helpers.CreateConnection(m.IP)
					encodedmsg2 := helpers.CreateMessage(3, "/"+args[4]+"/"+args[3]+"/"+args[1]+":"+args[2]+"+"+msg)
					stat2 := network.Send(conn2, encodedmsg2)
					if stat2 {
						fmt.Println("msg sent")
					}

				} else if ip != "NULL" {
					conn := helpers.CreateConnection(ip + ":" + port)
					encodedmsg := helpers.CreateMessage(op, "/"+args[4]+"/"+args[3]+"/"+args[1]+":"+args[2]+"+"+msg)
					stat := network.Send(conn, encodedmsg)
					if stat {
						fmt.Println("msg sent")
					}
				}

			}
		}
	}

}
