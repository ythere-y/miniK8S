package controllerManager

import (
	"fmt"
	"minik8s/registry/node"
	"minik8s/registry/pod"
	"minik8s/registry/service"
)

var MemNodes []node.Node
var MemPods []pod.Pod
var MemServices []service.Service
var relations Relation

func Main() {
	fmt.Println("[Controller Manager] Main started!")

}

func memInit() {
	relations.PodstoNodeRela = make(map[string]string)
	relations.NodetoPodRela = make(map[string][]string)
	relations.ServicetoPodRela = make(map[string][]string)
}

func ControllerStartUp() {
	memInit()
	nodeControllerWatch()
	podControllerWatch()
	serviceControllerWatch()
	replicasetControllerWatch()
	RsWatch()
}
