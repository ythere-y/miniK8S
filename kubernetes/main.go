package kubernetes

import (
	"fmt"
	"io/ioutil"
	"minik8s/apiserver"
	"minik8s/config"
	"minik8s/constant"
	"minik8s/controllerManager"
	"minik8s/kubelet"
	"minik8s/registry/node"
	"time"
)

func Main() {
	//Test func
	fmt.Println("kubernetes start up~!")

	StartUpMaster()

}

//CreateMasterNode
/*
填充初始的master节点的信息
*/
func CreateMasterNode() {
	var nodest node.Node
	nodest.Name = config.Configs.MasterNodeName
	controllerManager.AddNode(nodest)
	//apiserver.CmdCreateNode(nodest)
	go kubelet.Main()
}

//StartUpMaster
/*
开启一个master节点
*/
func StartUpMaster() {
	controllerManager.Main()

	apiserver.Main()

	nodeFile, err := ioutil.ReadFile(constant.MasterNodeFile)
	if err != nil {
		panic(err)
	}
	//apiserver.CmdCreateNode(nodeFile)

	//CreateMasterNode()

	time.Sleep(2 * time.Second)

}
