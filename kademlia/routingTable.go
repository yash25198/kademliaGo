package kademlia

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

type RoutingTableInterface interface {
	AddNode()
	FindNode() (string, string, bool)
	GetBucketCount() int
}

type RoutingTable struct {
	Buckets     []Bucket
	BucketCount int
	NodeID      string
}

func (x RoutingTable) GetBucketCount() int {
	return x.BucketCount
}

func sigBit(x string) int {
	n := len(x)
	if x == "0" {
		return 160
	}

	fmt.Println("bucket:", 160-n)
	return 160 - n
}

func findXOR(x string, y string) big.Int {
	xhex, _ := hex.DecodeString(x)
	yhex, _ := hex.DecodeString(y)

	xbig := new(big.Int).SetBytes(xhex)
	ybig := new(big.Int).SetBytes(yhex)

	xorxy := new(big.Int).Xor(xbig, ybig)
	return *xorxy
}

func findXORSigBit(x string, y string) int {
	xhex, _ := hex.DecodeString(x)
	yhex, _ := hex.DecodeString(y)

	xbig := new(big.Int).SetBytes(xhex)
	ybig := new(big.Int).SetBytes(yhex)

	xorxy := new(big.Int).Xor(xbig, ybig)
	fmt.Println(xorxy.Text(2))
	return sigBit(xorxy.Text(2))

}

func (x *RoutingTable) AddNode(newnode Node) {
	i := findXORSigBit(x.NodeID, newnode.NodeID)
	if i == 160 {
		return
	}
	fmt.Println(1, i, x.Buckets[i].GetCountList())
	x.Buckets[i].AppendToBucket(newnode)
}

func (x *RoutingTable) FindNode(nodeID string, requester string) (string, string, string, bool) {
	fmt.Println(x.NodeID, nodeID)
	i := findXORSigBit(x.NodeID, nodeID)
	//     fmt.Println("here:", nodeID, found, i, p, x.Buckets[p].GetNodeFromList(0))
	// if found {
	// 	return x.NodeID, "192.168.0.85", "8282", true
	// }
	p := i
	low := big.NewInt(2)
	e := big.NewInt(160)
	one := big.NewInt(1)
	low.Exp(low, e, nil).Sub(low, one)
	var nearest Node = Node{
		NodeID: "NULL",
		IP:     "NULL",
		Port:   "NULL",
	}
	fmt.Println("here:", nodeID, i, "what", x.Buckets[p].GetCountList())
	// for p < 160 {
	// 	fmt.Println("heeeerre")
	// 	if x.Buckets[p].GetCountList() > 1 || (x.Buckets[p].GetCountList() == 1 && x.Buckets[p].GetNodeFromList(0).NodeID != requester) {
	// 		fmt.Println(p)
	// 		break
	// 	} else {
	// 		p++
	// 	}

	// }
	// fmt.Println("heeeerre")
	// if p == 160 {
	// 	p = i
	// 	for p >= 0 {
	// 		if x.Buckets[p].GetCountList() > 0 {
	// 			fmt.Println(p)
	// 			break
	// 		} else {
	// 			p--
	// 		}
	// 		fmt.Println(p)

	// 	}
	// }
	// if p < 0 {
	// 	fmt.Println("herelol")
	// 	return "NULL", "NULL", "NULL", false
	// }
	// fmt.Println("here2:", nodeID)
	size := x.Buckets[p].GetCountList()
	fmt.Println("here3size:", size, p, x.Buckets[p])
	mapping := x.Buckets[p].GetMap()
	_, found := mapping[requester]
	if x.Buckets[p].GetCountList() == 1 && !found {
		for k := range mapping {
			ele := mapping[k]
			fmt.Println("First Element with loop", ele)
			return ele.NodeID, ele.IP, ele.Port, ele.NodeID == nodeID
		}
	}
	if size > 1 {
		for k := range mapping {
			tempNode := mapping[k]
			fmt.Println(tempNode, low)
			tempLow := findXOR(tempNode.NodeID, nodeID)
			if tempLow.Cmp(big.NewInt(0)) == 0 {
				fmt.Println("oof here", tempNode)
				return tempNode.NodeID, tempNode.IP, tempNode.Port, true
			}
			if tempLow.Cmp(low) < 0 && tempNode.NodeID != requester {
				low = &tempLow
				nearest = tempNode
			}
		}
		fmt.Println("damn here:", nearest)
		return nearest.NodeID, nearest.IP, nearest.Port, false
	}
	for j := 0; j < 160; j++ {
		if j == p {
			continue
		}
		sizej := x.Buckets[j].GetCountList()
		fmt.Println(j, sizej)
		mappingj := x.Buckets[j].GetMap()
		if sizej > 0 {
			for k := range mappingj {
				tempNode := mappingj[k]
				fmt.Println(tempNode)
				tempLow := findXOR(tempNode.NodeID, nodeID)
				// if tempLow.Cmp(big.NewInt(0)) == 0 {
				// 	fmt.Println("oof here", tempNode)
				// 	return tempNode.NodeID, tempNode.IP, tempNode.Port, true
				// }
				if tempLow.Cmp(low) < 0 && tempNode.NodeID != requester {
					low = &tempLow
					nearest = tempNode
				}
			}
		}
	}
	fmt.Println("nearest here :", nearest)
	return nearest.NodeID, nearest.IP, nearest.Port, false
}

func CreatRoutingTable(count int, nodeID string) RoutingTable {
	routingTable := RoutingTable{}
	routingTable.Buckets = make([]Bucket, count)
	for i := 0; i < count; i++ {
		routingTable.Buckets[i] = CreateBucket()
	}
	routingTable.BucketCount = count
	routingTable.NodeID = nodeID
	return routingTable
}
