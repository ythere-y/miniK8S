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
)

//StartUpMaster
/*
开启一个master节点
*/
func StartUpMaster() {
	var err error
	err = environment.EtcdStartUp(config.Configs.MasterIP)
	if err != nil {
		panic(err)
		return
	}
	err = environment.FlannelStartUp(config.Configs.EtcdIp)
	if err != nil {
		panic(err)
		return
	}
	controllerManager.ControllerStartUp()

	apiserver.Main()

	nodeFile, err := ioutil.ReadFile(constant.MasterNodeFile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("read file:\n%v\n", string(nodeFile))
	controllerManager.CreateMasterNode(nodeFile)

	kubelet.StartUp()
	kubelet.IptablesInit()
}
