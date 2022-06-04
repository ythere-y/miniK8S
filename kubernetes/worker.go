package kubernetes

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"io/ioutil"
	"minik8s/apiserver"
	"minik8s/config"
	"minik8s/constant"
	"minik8s/kubelet"
	"minik8s/lab/etcd"
	"minik8s/registry/node"
	"time"
)

var tmpNode node.Node

func StartUpWorker() {
	var err error
	if err != nil {
		panic(err)
		return
	}

	nodeFile, err := ioutil.ReadFile(constant.MasterNodeFile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("read file:\n%v\n", string(nodeFile))
	//err = environment.FlannelStartUp(config.Configs.MasterIP)
	kubelet.StartUp()
	kubelet.IptablesInit()
}

func JoinAsWorker(args []string) {
	var err error

	config.Configs.EtcdIp = args[0]

	apiserver.Main()

	nodeFile, err := ioutil.ReadFile(args[1])
	if err != nil {
		panic(err)
		return
	}
	fmt.Printf("read file:\n%v\n", string(nodeFile))
	cur := time.Now()
	apiserver.CmdCreateNode(nodeFile, cur)
	tmpNode = node.NodeYamlToNode(node.ParseNodeYaml(nodeFile))
	getstr, _ := json.Marshal(tmpNode)
	fmt.Printf("tmpnode get :\n%v\n", string(getstr))

	config.SetConfigThisNode(tmpNode)

	CreateNodeReplayWatch(cur)

}

func createNodeFail() {
	fmt.Println("join failed~!")
}
func CreateNodeReplayWatch(cur time.Time) {
	watchName := etcd.SetKey(
		etcd.SetPrefix(constant.ReplayPrefix),
		etcd.AppendName(cur.String()))
	apiserver.SyncWatchWithTime(watchName, handleNodeCreateReply, createNodeFail, constant.WaitReplyTime)
}
func handleNodeCreateReply(event *clientv3.Event) error {
	var err error
	err = nil

	switch event.Type {
	case mvccpb.PUT:
		if string(event.Kv.Value) == constant.ReplayOK {
			fmt.Println("get OK replay, start up worker")
			StartUpWorker()
		} else if string(event.Kv.Value) == constant.ReplayERROR {
			fmt.Println("get ERROR replay, end up")
			createNodeFail()
		}
	}
	return err
}
