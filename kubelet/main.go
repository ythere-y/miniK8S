package kubelet

import (
	"fmt"
	"minik8s/apiserver"
	"minik8s/constant"
	. "minik8s/lab/etcd"
	"minik8s/utils"
)

func CreateKubelet(listen string) {
	apiserver.SyncWatch(listen, kubelethandler)
}

func Main() {
	fmt.Println("hello world! Kubelet Main started!")

	watchName := SetKey(
		SetPrefix(constant.RegistryPrefix),
		SetSourceType(constant.NodeSourceName),
		SetNodeName("node_1"),
	)

	CreateKubelet(watchName)

	utils.HoldPro()
}
