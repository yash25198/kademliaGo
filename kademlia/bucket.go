package kademlia

// import (
// 	"fmt"
// )

type Node struct {
	IP     string
	Port   string
	NodeID string
}
type bucket struct {
	List      map[string]Node
	WaitList  map[string]Node
	Active    map[string]bool
	CountList int
	CountWait int
	MaxCount  int
}

type Bucket interface {
	AppendToBucket(newNode Node)
	RemoveFromHead(nodeID string)
	GetActive(nodeid string) bool
	GetMap() map[string]Node
	GetCountList() int
	GetCountWait() int
	GetMaxCount() int
}

func (x bucket) GetActive(nodeid string) bool {
	return x.Active[nodeid]
}
func (x bucket) GetMap() map[string]Node {
	return x.List
}

func (x bucket) GetCountList() int {
	return x.CountList
}

func (x bucket) GetCountWait() int {
	return x.CountWait
}
func (x bucket) GetMaxCount() int {
	return x.MaxCount
}

func (x *bucket) AppendToBucket(newNode Node) {
	if x.Active[newNode.NodeID] {
		x.List[newNode.NodeID] = newNode
		return
	}
	if x.CountList == x.MaxCount {
		x.WaitList[newNode.NodeID] = newNode
		x.CountWait += 1
		return
	}
	x.List[newNode.NodeID] = newNode
	x.CountList += 1
	x.Active[newNode.NodeID] = true
	// fmt.Println("list", x, x.List, x.CountList)
}

func (x *bucket) RemoveFromHead(nodeID string) {
	delete(x.List, nodeID)
	x.CountList -= 1
}

func CreateBucket() Bucket {
	return &bucket{
		List:     map[string]Node{},
		WaitList: map[string]Node{},
		Active:   map[string]bool{},
		MaxCount: 20,
	}
}
