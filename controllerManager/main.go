package controllerManager

import (
	"fmt"
	"minik8s/registry/node"
	"minik8s/registry/pod"
	"minik8s/registry/service"
	"minik8s/utils"
)

var MemNodes []node.Node
var MemPods []pod.Pod
var MemServices []service.Service
var relations Relation

func Main() {
	fmt.Println("[Controller Manager] Main started!")

	utils.HoldPro()
}

func init() {
	relations.PodstoNodeRela = make(map[string]string)
	relations.NodetoPodRela = make(map[string][]string)
	relations.ServicetoPodRela = make(map[string][]string)
}
