package kubelet

import (
	"fmt"
	. "minik8s/lab/etcd"
	"minik8s/utils"
	"time"
)

func CreateKubelet(listen string) {
	SyncWatch(listen, kubelethandler)
}

func reqeustPutTest() {
	SyncPut(SetKey(
		SetPrefix("registry"),
		SetSourceType("pods"),
		SetNameSpace("default"),
		SetNodeName("node1"),
		SetPodName("pod1")), "hello")

}

func KubeletMain() {
	fmt.Println("hello world! Kubelet Main started!")

	watchName := SetKey(
		SetPrefix("registry"),
		SetSourceType("pods"),
		SetNameSpace("default"),
		SetNodeName("node1"),
	)

	CreateKubelet(watchName)

	// sleep
	time.Sleep(3 * time.Second)

	reqeustPutTest()

	utils.HoldPro()
}
