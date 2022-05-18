package kubelet

import (
	"fmt"
	. "minik8s/lab/etcd"
	"minik8s/utils"
)

func CreateKubelet(listen string) {
	SyncWatch(listen, kubelethandler)
}

func KubeletMain() {
	fmt.Println("hello world! Kubelet Main started!")

	watchName := SetKey(
		SetPrefix("registry"),
		SetSourceType("pods"),
		SetNameSpace("default"),
		SetNodeName("node_1"),
	)

	CreateKubelet(watchName)

	utils.HoldPro()
}
