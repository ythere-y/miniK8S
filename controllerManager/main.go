package controllerManager

import (
	"fmt"
	"minik8s/constant"
	. "minik8s/lab/etcd"
	"minik8s/utils"
	"time"
)

func CreateControllerManager(listen string) {
	SyncWatch(listen, controllerManagerHandler)
}

func reqeustPutTest() {
	// sleep
	time.Sleep(3 * time.Second)
	SyncPut(SetKey(
		SetPrefix("controller"),
		SetSourceType("pods"),
		SetNameSpace("id_1")),
		"pod_yaml")

}

func Main() {
	fmt.Println("Controller Manager Main started!")

	// pod watcher
	watchName := SetKey(
		SetPrefix(constant.ControllerPrefix),
		SetSourceType("pods"),
	)
	CreateControllerManager(watchName)

	// replicaset watcher
	rsWatchName := SetKey(
		SetPrefix(constant.ControllerPrefix),
		SetSourceType("replicaset"),
	)
	CreateControllerManager(rsWatchName)

	utils.HoldPro()
}
