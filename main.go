package main

import "minik8s/lab/etcd"

// func main() {
// 	//service.SerMain()
// 	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
// 	if err != nil {
// 		panic(err)
// 	}
// 	source := dksdk.Resource{
// 		CPUShares: 2,
// 		Memory:    128000000,
// 	}
// 	id := dksdk.CreateContainer(cli, "nginx", nil, source, "testCreate", nil, nil, "14c453757e37")
// 	time.Sleep(time.Second * 1)
// 	dksdk.StartContainer(id, cli)
// 	time.Sleep(time.Second * 1)
// 	flag := dksdk.IsRun(cli, id)
// 	if flag {
// 		fmt.Printf("container %s is running\n", id)
// 	}
// 	dksdk.ContainerStat(cli, id)
// 	dksdk.ListContainer(cli)
// 	//time.Sleep(time.Second * 1)
// 	//dksdk.StopContainer(id, cli)
// 	//time.Sleep(time.Second * 3)
// 	//id, err = dksdk.RemoveContainer(id, cli)
// 	//if err == nil {
// 	//	fmt.Println("删除容器", id, "成功")
// 	//}
// 	// dksdk.CreateContainerInBackground()
// }

var rootName string

func main() {
	etcd.LabMain() //关于etcd的测试
	//cmd.RootCmdRun() // 关于命令行的测试
	//service.OutPutFmtTest() //关于格式化输出的测试
	//service.SerReadTest() // 关于读取yamle文件建立service的测试

	//circle.CircleTest() // 关于循环import的测试
	return

}
