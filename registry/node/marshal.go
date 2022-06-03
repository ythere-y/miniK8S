package node

import (
	"gopkg.in/yaml.v2"
	"time"
)

func ParseNodeYaml(filecontext []byte) NodeYaml {
	var newNodeYaml NodeYaml
	err := yaml.Unmarshal(filecontext, &newNodeYaml)
	if err != nil {
		panic(err)
	}
	return newNodeYaml
}
func NodeYamlToNode(nodeB NodeYaml) Node {
	var retNode Node
	retNode.Addr = nodeB.Addr
	retNode.Name = nodeB.Name
	retNode.Capacity = 16
	retNode.Status = NODE_OK
	retNode.CreatTime = time.Now()
	return retNode
}
