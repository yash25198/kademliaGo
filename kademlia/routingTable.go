package kademlia

import (
	"encoding/hex"
	// "fmt"
	"math/big"
	// "strings"
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

func sigBit(x string) (int, bool) {
	n := len(x)
	if x == "0" {
		return 160, true
	}

	// fmt.Println("bucket:", 160-n)
	return 160 - n, false
}

func findXORSigBit(x string, y string) (int, bool) {
	xhex, _ := hex.DecodeString(x)
	yhex, _ := hex.DecodeString(y)

	xbig := new(big.Int).SetBytes(xhex)
	ybig := new(big.Int).SetBytes(yhex)

	xorxy := new(big.Int).Xor(xbig, ybig)
	// fmt.Println(xorxy.Text(2))
	return sigBit(xorxy.Text(2))

}

func (x *RoutingTable) AddNode(newnode Node) {
	i, found := findXORSigBit(x.NodeID, newnode.NodeID)
	if found {
		return
	}
	// fmt.Println(1, i, x.Buckets[i].GetCountList())
	x.Buckets[i].AppendToBucket(newnode)
}

func (x *RoutingTable) FindNode(nodeID string) (string, string, string, bool) {
	// fmt.Println(x.NodeID, nodeID)
	i, found := findXORSigBit(x.NodeID, nodeID)
	// fmt.Println("here:", nodeID, found, i, p, x.Buckets[p].GetNodeFromList(0))
	if found {
		return x.NodeID, "192.168.0.85", "8282", true
	}
	p := i
	low := -1
	var nearest Node
	// fmt.Println("here:", nodeID, i, "what", x.Buckets[p].GetCountList())
	for p < 160 {
		// fmt.Println("heeeerre")
		if x.Buckets[p].GetCountList() > 0 {
			// fmt.Println(p)
			break
		} else {
			p++
		}

	}
	// fmt.Println("heeeerre")
	if p == 160 {
		p = i
		for p >= 0 {
			if x.Buckets[p].GetCountList() > 0 {
				break
			} else {
				p--
			}

		}
	}
	if p < 0 {
		// fmt.Println("herelol")
		return "NULL", "NULL", "NULL", false
	}
	// fmt.Println("here2:", nodeID)
	size := x.Buckets[i].GetCountList()
	// fmt.Println("here3size:", size)
	for j := 0; j < size; j++ {
		tempNode := x.Buckets[p].GetNodeFromList(j)
		// fmt.Println(tempNode)
		tempLow, found := findXORSigBit(tempNode.NodeID, nodeID)
		if found {
			// fmt.Println("oof here", tempNode)
			return tempNode.NodeID, tempNode.IP, tempNode.Port, true
		}
		if tempLow > low {
			low = tempLow
			nearest = tempNode
		}
	}
	// fmt.Println("damn here:", nearest)
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
