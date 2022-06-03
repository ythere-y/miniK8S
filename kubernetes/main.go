package kubernetes

import (
	"fmt"
	"io/ioutil"
	"minik8s/apiserver"
	"minik8s/constant"
	"minik8s/constant/config"
	"minik8s/controllerManager"
	"minik8s/kubelet"
	"minik8s/registry/node"
	"os"
	"os/exec"
	"time"
)

func Main() {
	//Test func
	fmt.Println("hello world")
	//export NODE=10.119.11.72
	//docker run -p 2379:2379 quay.io/coreos/etcd:v3.3.1 /usr/local/bin/etcd --advertise-client-urls http://${NODE}:2379 --listen-client-urls http://0.0.0.0:2379
	var (
		cmdStr string
		cmd    *exec.Cmd
	)

	cmdStr = "ls"
	cmd = exec.Command(cmdStr)

	if err := cmd.Run(); err != nil {
		panic(err)
	}

	cmdStr = "export"
	cmd = exec.Command(cmdStr, "NODE=localhost")

	if err := cmd.Run(); err != nil {
		panic(err)
	}

	//fmt.Printf("command = %s, result = %s\n", cmdStr, string(output))

	os.Exit(1)

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
	go controllerManager.Main()

	apiserver.Main()

	nodeFile, err := ioutil.ReadFile(constant.MasterNodeFile)
	if err != nil {
		panic(err)
	}
	apiserver.CmdCreateNode(nodeFile)

	CreateMasterNode()

	time.Sleep(2 * time.Second)

}
