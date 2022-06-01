package kubernetes

import (
	"fmt"
	"minik8s/apiserver"
	"minik8s/config"
	"minik8s/controllerManager"
	"minik8s/kubelet"
	"minik8s/node"
	"minik8s/scheduler"
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
