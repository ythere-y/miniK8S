package kubelet

import (
	"fmt"
	"minik8s/apiserver"
	"minik8s/config"
	"minik8s/constant"
	. "minik8s/lab/etcd"
	"minik8s/utils"
)

var NodeName string = ""

func CreateKubelet(listen string) {
	apiserver.SyncWatch(listen, kubelethandler)
}

func Main() {
	fmt.Println("[Kubelet] [name = " + NodeName + "] Main started!")

	watchName := SetKey(
		SetPrefix(constant.RelationPrefix),
		SetSourceType(constant.NodeSourceName),
		SetNodeName(NodeName),
	)

	CreateKubelet(watchName)

	utils.HoldPro()
}

func StartUp() {
	fmt.Printf("[Kubelete] [name = %v ] start up \n")
}

func init() {
	NodeName = config.Configs.MasterNodeName
}
