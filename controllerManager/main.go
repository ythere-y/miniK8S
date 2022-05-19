package controllerManager

import (
	"fmt"
	"minik8s/apiserver"
	"minik8s/constant"
	. "minik8s/lab/etcd"
	"minik8s/utils"
	"time"
)

func CreateControllerManager() {
	var watchName string
	watchName = SetKey(
		SetPrefix(constant.ControllerPrefix),
		JustAppend(constant.CREATE),
		SetSourceType(constant.PodSourceName),
	)
	apiserver.SyncWatch(watchName, createPod)

	watchName = SetKey(
		SetPrefix(constant.ControllerPrefix),
		JustAppend(constant.DELETE),
		SetSourceType(constant.PodSourceName),
	)
	apiserver.SyncWatch(watchName, deletePod)

}

func reqeustPutTest() {
	// sleep
	time.Sleep(3 * time.Second)
	apiserver.SyncPut(SetKey(
		SetPrefix(constant.ControllerPrefix),
		SetSourceType("pods"),
		SetNameSpace("id_1")),
		"pod_yaml")

}

func Main() {
	fmt.Println("Controller Manager Main started!")

	CreateControllerManager()

	utils.HoldPro()
}
