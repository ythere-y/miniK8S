package kubelet

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"minik8s/src/pod"
	"minik8s/utils"
)

func kubelethandler(event *clientv3.Event) error {
	var err error
	fmt.Printf("kubelet handling key = %v, value = %v\n", string(event.Kv.Key), string(event.Kv.Value))
	switch event.Type {
	case mvccpb.DELETE:
		// 删除一个pod的操作
	case mvccpb.PUT:
		// 增加/修改 一个pod的操作
		var podInfo pod.Pod
		err = json.Unmarshal(event.Kv.Value, &podInfo)
		utils.HandleError("unmarshal pod error", err)
		pod.CreateAndRunPod(podInfo)
	}
	return err
}
