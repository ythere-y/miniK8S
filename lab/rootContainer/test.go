package rootContainer

import (
	"github.com/docker/docker/client"
	"minik8s/lab/dksdk"
)

func Test() {
	id0 := dksdk.RunRootContainer("bb")
	print(id0)
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
	}
	source := dksdk.Resource{
		CPUShares: 2,
		Memory:    128000000,
	}

	//binds := []string{"D:\\Schoolwork\\2022_spring\\CloudComputing\\labs\\Minik8s\\data:/data"}
	id := dksdk.CreateContainer(cli, "nginx", nil, source, "name1", nil, nil, "", id0)
	dksdk.StartContainer(id, cli)
}
