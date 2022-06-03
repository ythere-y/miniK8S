package node

import (
	"fmt"
	"minik8s/utils"
	"strconv"
	"time"
)

type NodeStatus string

const (
	NODE_OK   NodeStatus = "OK"
	NODE_FAIL NodeStatus = "FAIL"
)

type Node struct {
	Capacity  uint32
	Addr      string
	Name      string
	Status    NodeStatus
	CreatTime time.Time
}

type NodeBasic struct {
	Addr string
	Name string
}

var blockSize = 20

func NodeBasicToNode(nodeB NodeBasic) Node {
	var retNode Node
	retNode.Addr = nodeB.Addr
	retNode.Name = nodeB.Name
	retNode.Capacity = 16
	retNode.Status = NODE_OK
	retNode.CreatTime = time.Now()
	return retNode
}

func NodePreDisplay() {
	fmt.Printf("%-"+strconv.Itoa(blockSize)+"s", "NAME")
	fmt.Printf("%-"+strconv.Itoa(blockSize)+"s", "STATUS")
	fmt.Printf("%-"+strconv.Itoa(blockSize)+"s", "ADDR")
	fmt.Printf("%-"+strconv.Itoa(blockSize)+"s", "AGE")
	fmt.Println()
}

func (node Node) Display() {
	fmt.Printf("%-"+strconv.Itoa(blockSize)+"s", node.Name)
	fmt.Printf("%-"+strconv.Itoa(blockSize)+"v", node.Status)
	fmt.Printf("%-"+strconv.Itoa(blockSize)+"v", node.Addr)
	fmt.Printf("%-"+strconv.Itoa(blockSize)+"v", utils.GetAge(node.CreatTime))
	fmt.Println()
}
