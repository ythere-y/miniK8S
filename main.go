package main

import (
	"minik8s/lab/dksdk"
	"github.com/docker/docker/client"
	"fmt"
	"time"
)

func main() {
	//service.SerMain()
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
	}
	id := dksdk.CreateContainer(cli, "library/alpine", []string{"echo", "hello world"}, "testCreate", nil, nil)
	fmt.Printf("%s\n", id)
	time.Sleep(time.Second * 1)
	dksdk.StartContainer(id, cli)
	time.Sleep(time.Second * 2)
	dksdk.ListContainer(cli)
	time.Sleep(time.Second * 1)
	dksdk.StopContainer(id, cli)
	time.Sleep(time.Second * 3)
	id, err = dksdk.RemoveContainer(id, cli)
	if err == nil {
		fmt.Println("删除容器", id, "成功")
	}
}
