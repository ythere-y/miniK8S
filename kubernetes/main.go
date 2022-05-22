package kubernetes

import (
	"minik8s/apiserver"
	"minik8s/config"
	"minik8s/controllerManager"
	"minik8s/kubelet"
	"minik8s/node"
	"minik8s/scheduler"
	"time"
)

func Main() {
	//Test func
}

//CreateMasterNode
/*
填充初始的master节点的信息
*/
func CreateMasterNode() {
	var nodest node.NodeStatus
	nodest.Name = config.Configs.MasterNodeName
	controllerManager.AddNode(nodest)
	apiserver.CreateNode(nodest)
	go kubelet.Main()
}

//StartUpMaster
/*
开启一个master节点
*/
func StartUpMaster() {
	go scheduler.Main()
	go controllerManager.Main()
	go apiserver.Main()

	CreateMasterNode()

	time.Sleep(2 * time.Second)

}
