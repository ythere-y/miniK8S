package service

import (
	"fmt"
	"github.com/docker/docker/client"
	"minik8s/lab/dksdk"
	"time"
)

func Dksdktest() {

	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
	}
	id := dksdk.CreateContainer(cli, "library/nginx", nil, "testCreate", nil, nil)
	time.Sleep(time.Second * 1)
	dksdk.StartContainer(id, cli)
	time.Sleep(time.Second * 1)
	flag := dksdk.IsRun(cli, id)
	if flag {
		fmt.Printf("container %s is running\n", id)
	}
	dksdk.ContainerStat(cli, id)
	dksdk.ListContainer(cli)
	time.Sleep(time.Second * 1)
	dksdk.StopContainer(id, cli)
	time.Sleep(time.Second * 3)
	id, err = dksdk.RemoveContainer(id, cli)
	if err == nil {
		fmt.Println("删除容器", id, "成功")
	}
}
