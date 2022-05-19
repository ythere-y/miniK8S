package controllerManager

import (
	"fmt"
	"minik8s/apiserver"
	"minik8s/constant"
	. "minik8s/lab/etcd"
	"minik8s/utils"
	"time"
)

func CreateControllerManager(listen string) {
	apiserver.SyncWatch(listen, controllerManagerHandler)
}

func reqeustPutTest() {
	// sleep
	time.Sleep(3 * time.Second)
	apiserver.SyncPut(SetKey(
		SetPrefix("controller"),
		SetSourceType("pods"),
		SetNameSpace("id_1")),
		"pod_yaml")

}

func Main() {
	fmt.Println("Controller Manager Main started!")

	watchName := SetKey(
		SetPrefix(constant.ControllerPrefix),
		SetSourceType("pods"),
	)

	CreateControllerManager(watchName)

	utils.HoldPro()
}
