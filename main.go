package main

import (
	"fmt"
	"github.com/docker/docker/client"
	"minik8s/lab/circle"
	"minik8s/lab/dksdk"
	"time"
)

var rootName string

func main() {
	//cmd.RootCmdRun()// 关于命令行的测试

	//service.SerReadTest()	// 关于读取yamle文件建立service的测试

	circle.CircleTest() // 关于循环import的测试
	return

	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
	}
	id := dksdk.CreateContainer(cli, "library/alpine", []string{"echo", "hello world"}, "testCreate", nil, nil)
	fmt.Printf("%s\n", id)
	time.Sleep(time.Second * 1)
	dksdk.StartContainer(id, cli)
	time.Sleep(time.Second * 40)
	dksdk.StopContainer(id, cli)
	time.Sleep(time.Second * 3)
	id, err = dksdk.RemoveContainer(id, cli)
	if err == nil {
		fmt.Println("删除容器", id, "成功")
	}
}
