package controllerManager

import (
	"fmt"
	"minik8s/apiserver"
	"minik8s/constant"
	. "minik8s/lab/etcd"
	"minik8s/registry/node"
	"minik8s/utils"
	"time"
)

var KNodes []node.Node

func AddNode(status node.Node) {
	KNodes = append(KNodes, status)

}

func CreateControllerManager() {
	var watchName string
	watchName = SetKey(
		SetPrefix(constant.ControllerPrefix),
		SetSourceType(constant.PodSourceName),
		JustAppend(constant.CREATE),
	)
	apiserver.SyncWatch(watchName, createPod)

	watchName = SetKey(
		SetPrefix(constant.ControllerPrefix),
		SetSourceType(constant.PodSourceName),
		JustAppend(constant.DELETE),
	)
	apiserver.SyncWatch(watchName, deletePod)

	watchName = SetKey(
		SetPrefix(constant.ControllerPrefix),
		SetSourceType(constant.PodSourceName),
		JustAppend(constant.STOP),
	)
	apiserver.SyncWatch(watchName, stopPod)

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
	fmt.Println("[Controller Manager] Main started!")

	CreateControllerManager()

	utils.HoldPro()
}
