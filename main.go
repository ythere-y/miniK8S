package main

import (
	"minik8s/cmd"
	"minik8s/controllerManager"
	"minik8s/kubelet"
	"minik8s/scheduler"
	"minik8s/utils"
	"time"
)

var rootName string

func main() {
	//go environment.Main() //测试启动前的环境准备
	go controllerManager.Main() // 启动controller manager
	go scheduler.Main()         // 启动scheduler
	go kubelet.Main()           // 启动kubelet
	//apiserver.ApiServerTest() //关于apiserver的测试
	//go etcd.LabMain() //关于etcd的测试

	time.Sleep(2 * time.Second)

	cmd.RootCmdRun() // 关于命令行的测试
	//service.OutPutFmtTest() //关于格式化输出的测试
	//service.SerReadTest() // 关于读取yamle文件建立service的测试

	//circle.CircleTest() // 关于循环import的测试
	utils.HoldPro() // 阻塞进程防止运行结束
	return

}
