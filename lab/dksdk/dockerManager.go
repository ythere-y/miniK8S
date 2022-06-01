package dksdk

// SDK官方指南https://docs.docker.com/engine/api/sdk/examples/

import (
	"bytes"
	"context"
	"fmt"

	//"github.com/docker/docker/pkg/stdcopy"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func RunRootContainer(image string, cmd []string, resource Resource, name string, volumn map[string]struct{}) string {
	return ""
}

func CreateContainer(cli *client.Client, image string, cmd []string, resource Resource, name string, binds []string, exports nat.PortSet, hostport string, network string) string {

	ctx := context.Background()
	reader, err := cli.ImagePull(ctx, "docker.io/"+image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resources := container.Resources{
		CPUShares: resource.CPUShares,
		Memory:    resource.Memory,
	}
	hostconfig := &container.HostConfig{
		Resources: resources,
		Binds:     binds,
	}
	if network != "" {
		hostconfig.NetworkMode = container.NetworkMode("container:" + network)
	}

	if hostport != "" {
		hostconfig.PortBindings = nat.PortMap{
			nat.Port(fmt.Sprintf("80/tcp")): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: hostport,
				},
			},
		}
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   cmd,
		Tty:   true,
		ExposedPorts: nat.PortSet{
			"80/tcp": {},
		},
		//Shell: []string{"cmd.exe", "/c", "-P"},
	}, hostconfig, nil, nil, name)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ID: %s\n", resp.ID)
	return resp.ID
}

// 启动
func StartContainer(containerID string, cli *client.Client) {
	err := cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if err == nil {
		fmt.Println("容器", containerID, "启动成功")
	} else {
		panic(err)
		//fmt.Println("容器", containerID, "启动失败")
	}
}

// 停止
func StopContainer(containerID string, cli *client.Client) {
	timeout := time.Second * 10
	err := cli.ContainerStop(context.Background(), containerID, &timeout)
	if err != nil {
		fmt.Println("容器", containerID, "停止失败")
	} else {
		fmt.Printf("容器%s已经被停止\n", containerID)
	}
}

// 删除
func RemoveContainer(containerID string, cli *client.Client) (string, error) {
	err := cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{})
	//log(err)
	return containerID, err
}

// docker ps -a
func ListContainer(cli *client.Client) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)

	}
	fmt.Println("container.ID,\t\t\t\t\t\t\tcontainer.Names,    container.Created,   container.Status,    container.Ports")
	for _, container := range containers {
		fmt.Println(container.ID, container.Names, container.Created, container.Status, container.Ports)
	}
}

func IsRun(cli *client.Client, containerID string) bool {
	stat, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return false
	}
	if !stat.State.Running {
		return false
	}
	return true
}

func ContainerStat(cli *client.Client, containerID string) {
	ctx := context.Background()
	containerStats, err := cli.ContainerStats(ctx, containerID, false)
	if err != nil {
		panic(err)
	}
	/**
	ContainerStats的返回的结构如下 注意这个Body的类型是io.ReadCloser 好奇怪的类型 下面我们给他转成json
	type ContainerStats struct {
		Body   io.ReadCloser `json:"body"`
		OSType string        `json:"ostype"`
	}
	*/
	fmt.Println(containerStats)
	fmt.Println("containerStats.Body的内容是: ", containerStats.Body)
	buf := new(bytes.Buffer)
	//io.ReadCloser 转换成 Buffer 然后转换成json字符串
	buf.ReadFrom(containerStats.Body)
	newStr := buf.String()
	fmt.Printf(newStr)
}

// 后台运行容器，相当于键入 docker run -d bfirsh/reticulate-splines
func CreateContainerInBackground() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
	}

	imageName := "bfirsh/reticulate-splines"

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)

}

// 列出所有镜像
func ListImages() {
	//ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.ID)
	}
}

// 打印特定容器的日志
func LogsByContainerId() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
	}

	options := types.ContainerLogsOptions{ShowStdout: true}
	// Replace this ID with a container that really exists
	out, err := cli.ContainerLogs(ctx, "f1064a8a4c82", options)
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)
}

// 列出并管理容器
func ListAndManageContainer() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.ID)
		// 停止容器
		if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
			panic(err)
		}
		// 删除容器
		cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{})
	}
}
