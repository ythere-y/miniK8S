package scheduler

import (
	"minik8s/apiserver"
	"minik8s/utils"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func schedulerHandler(event *clientv3.Event) error {
	var err error
	var podName string
	podName = string(event.Kv.Value)

	// 对pod进行一些查询等操作，找到对应的node
	// TODO:需要完善scheduler的具体算法
	// 暂时先全部分配到node_1上
	err = apiserver.DistributePodtoNode("node_1", podName)
	utils.HandleError("Distribute pod to node error", err)

	return err
}
