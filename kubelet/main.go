package kubelet

import (
	"fmt"
	"minik8s/apiserver"
	"minik8s/config"
	"minik8s/registry/pod"
	"minik8s/utils"
)

var NodeName string = ""
var (
	KPods    []pod.Pod
	PodsInfo map[string]pod.Pod
)

func CreateKubelet(listen string) {
	apiserver.SyncWatch(listen, nodehandler)
}

func Main() {

	utils.HoldPro()
}

func StartUp() {
	fmt.Printf("[Kubelete] [name = %v ] [this IP = %v] start up \n", config.Configs.ThisNode.Name, config.Configs.ThisNode.Addr)
	NodeName = config.Configs.LocalNodeName
}

func init() {
	PodsInfo = make(map[string]pod.Pod)
}
