package kubernetes

import (
	"fmt"
	"io/ioutil"
	"minik8s/apiserver"
	"minik8s/config"
	"minik8s/constant"
	"minik8s/controllerManager"
	"minik8s/environment"
	"minik8s/kubelet"
	"minik8s/registry/node"
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
	var err error
	err = environment.EtcdStartUp()
	if err != nil {
		panic(err)
		return
	}
	err = environment.FlannelStartUp(config.Configs.EtcdIp)
	if err != nil {
		panic(err)
		return
	}
	controllerManager.Main()

	apiserver.Main()

	nodeFile, err := ioutil.ReadFile(constant.MasterNodeFile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("read file:\n%v\n", string(nodeFile))

	apiserver.CmdCreateNode(nodeFile)

}

func JoinAsWorker(args []string) {
	var err error
	err = environment.FlannelStartUp(args[0])
	if err != nil {
		panic(err)
		return
	}
	config.Configs.EtcdIp = args[0]
	kubelet.StartUp()
	apiserver.Main()

	nodeFile, err := ioutil.ReadFile(args[1])
	if err != nil {
		panic(err)
		return
	}
	fmt.Printf("read file:\n%v\n", string(nodeFile))

	apiserver.CmdCreateNode(nodeFile)
}
