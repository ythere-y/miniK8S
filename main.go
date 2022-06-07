package main

import (
	"fmt"
	"minik8s/apiserver"
	"minik8s/test"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var rootName string

func tmph(event *clientv3.Event) error {
	fmt.Printf("time out got value = %v\n", string(event.Kv.Value))
	return nil
}

func nor(event *clientv3.Event) error {
	fmt.Printf("normal got value = %v\n", string(event.Kv.Value))
	return nil
}

func fail() {
	fmt.Println("fail!")
}
func timeoutTest() {

	apiserver.SyncWatchWithTime("/test", tmph, fail, 5*time.Second)
}

func normalTest() {
	apiserver.SyncWatch("/test", nor)
}
func watchwithTimeTest() {

	go timeoutTest()
	go normalTest()
	time.Sleep(2 * time.Second)
	apiserver.SyncPut("/test", "test1")
	apiserver.SyncPut("/test", "test11")
	apiserver.SyncPut("/test", "test11111")

	time.Sleep(2 * time.Second)
	apiserver.SyncPut("/test", "test2")
}

func Debug() {

	//kubernetes.StartUpMaster()
	//
	//apiserver.DisplayAllPodsInfo()
}

func main() {
	// cmd.RootCmdRun() // 关于命令行的测试
	//rootContainer.Test()

	//fmt.Println(dns.RunCmd("ls"))
	//out := dns.ParseYaml("./dns/dns.yaml")
	//dns.C2Dock(out)

	//utils.CopyandRun()
	//slurm := utils.ParseYaml(utils.Dir + "test.yaml")
	//fmt.Print(slurm)
	//utils.Slu2File(slurm)
	//utils.Submit(utils.Dir+"test.yaml", utils.Dir+"matrix.cu")
	//utils.GetStat("matrix")
	//utils.GetRes("matrix")

	//rootContainer.Test()
	//utils.GetRes("test002")
	//fmt.Println(err.Error())

	//log.Println("hello world")
	test.FloatTest()
	//Debug()
	//shellScripts.Main()
	//rootContainer.Test()
	//watchwithTimeTest() // 关于定时watch的测试
	//controllerManager.ControllerStartUp()
	//kubernetes.Test() // kubernetes的测试部分
	//go kubernetes.Main() // kubernetes的测试部分
	//go cni.Main() // 测试CNI插件功能部分
	//go environment.Main() //测试启动前的环境准备
	//go controllerManager.Main() // 启动controller manager
	//go scheduler.Main()         // 启动scheduler
	//go kubelet.Main()           // 启动kubelet
	//apiserver.ApiServerTest() //关于apiserver的测试
	//go etcd.LabMain() //关于etcd的测试
	//config.Main()
	//time.Sleep(2 * time.Second)
	//go kubernetes.StartUpMaster()	// master节点的初始化startup
	//go cmd.RootCmdRun() // 关于命令行的测试
	//service.OutPutFmtTest() //关于格式化输出的测试
	//service.SerReadTest() // 关于读取yamle文件建立service的测试

	//circle.CircleTest() // 关于循环import的测试
	//utils.HoldPro() // 阻塞进程防止运行结束

	//time.Sleep(time.Second * 10)
	return

}
