package main

import (
	"fmt"
	"github.com/docker/docker/client"
	"minik8s/lab/dksdk"
	"minik8s/service"
	"time"
)

func main() {
	service.SerMain()
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
