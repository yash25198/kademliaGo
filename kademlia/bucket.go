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
	List      []Node
	WaitList  []Node
	Active    map[string]bool
	CountList int
	CountWait int
	MaxCount  int
}

type Bucket interface {
	AppendToBucket(newNode Node)
	RemoveFromHead()
	GetActive(nodeid string) bool
	GetNodeFromList(i int) Node
	GetCountList() int
	GetCountWait() int
	GetMaxCount() int
}

func (x bucket) GetActive(nodeid string) bool {
	return x.Active[nodeid]
}
func (x bucket) GetNodeFromList(i int) Node {
	return x.List[i]
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
		return
	}
	if x.CountList == x.MaxCount {
		x.WaitList = append(x.WaitList, newNode)
		x.CountWait += 1
		return
	}
	x.List = append(x.List, newNode)
	x.CountList += 1
	x.Active[newNode.NodeID] = true
	// fmt.Println("list", x, x.List, x.CountList)
}

func (x *bucket) RemoveFromHead() {
	x.List = x.List[1:]
	x.CountList -= 1
	if x.CountWait > 0 {
		x.AppendToBucket(x.WaitList[0])
		x.WaitList = x.WaitList[1:]
		x.CountWait -= 1
	}
}

func CreateBucket() Bucket {
	return &bucket{
		List:     []Node{},
		WaitList: []Node{},
		Active:   map[string]bool{},
		MaxCount: 20,
	}
}
