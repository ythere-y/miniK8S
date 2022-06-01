package localhost

import (
	"github.com/docker/docker/client"
	"minik8s/lab/dksdk"
)

//busybox run "wget -O- localhost" ==> nginx

func Test() {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
	}
	source := dksdk.Resource{
		CPUShares: 2,
		Memory:    128000000,
	}
	binds := []string{"D:\\Schoolwork\\2022_spring\\CloudComputing\\labs\\Minik8s\\data:/data"}
	id := dksdk.CreateContainer(cli, "busybox", nil, source, "name1", binds, nil, "", "")
	dksdk.StartContainer(id, cli)
	id2 := dksdk.CreateContainer(cli, "nginx", nil, source, "name2", binds, nil, "", id)
	dksdk.StartContainer(id2, cli)
}
