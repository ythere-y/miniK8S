package kubernetes

import (
	"fmt"
	"io/ioutil"
	"minik8s/apiserver"
	"minik8s/config"
	"minik8s/constant"
	"minik8s/controllerManager"
	"minik8s/environment"
	"time"
)

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
	cur := time.Now()
	apiserver.CmdCreateNode(nodeFile, cur)

}
