package scheduler

import (
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"minik8s/apiserver"
	"minik8s/src/pod"
	"minik8s/utils"
)

func schedulerHandler(event *clientv3.Event) error {
	var err error
	var podInfo pod.Pod
	err = json.Unmarshal(event.Kv.Value, &podInfo)
	utils.HandleError("Unmarshal pod error", err)
	// 对pod进行一些查询等操作，找到对应的node
	// 暂时先全部分配到node_1上
	err = apiserver.DistributePodtoNode("node_1", podInfo)
	utils.HandleError("Distribute pod to node error", err)

	return err
}
