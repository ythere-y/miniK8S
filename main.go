package main

import (
	"github.com/docker/docker/client"
	"minik8s/lab/dksdk"
)

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
	//cmd.RootCmdRun() // 关于命令行的测试

	//utils.CopyandRun()
	//slurm := utils.ParseYaml(utils.Dir + "test.yaml")
	//fmt.Print(slurm)
	//utils.Slu2File(slurm)
	//utils.Submit(utils.Dir+"test.yaml", utils.Dir+"test002.cu")
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
	}
	source := dksdk.Resource{
		CPUShares: 2,
		Memory:    128000000,
	}
	id := dksdk.CreateContainer(cli, "nginx", nil, source, "name", nil, nil, "8089", "")
	dksdk.StartContainer(id, cli)
	//fmt.Println(err.Error())

	//service.OutPutFmtTest() //关于格式化输出的测试
	//service.SerReadTest() // 关于读取yamle文件建立service的测试

	//circle.CircleTest() // 关于循环import的测试
	return

}
