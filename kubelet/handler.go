package kubelet

import (
	"encoding/json"
	"fmt"
	"minik8s/pod"

	"minik8s/apiserver"
	"minik8s/pod"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func kubelethandler(event *clientv3.Event) error {
	var err error
	fmt.Printf("kubelet handling key = %v, value = %v\n", string(event.Kv.Key), string(event.Kv.Value))
	switch event.Type {
	case mvccpb.DELETE:
		// 删除一个pod的操作
	case mvccpb.PUT:
		// 增加/修改 一个pod的操作
		pod_name := string(event.Kv.Value)
		var podInfo pod.Pod
		podInfo = apiserver.GetPodInfo(pod_name)
		//err = json.Unmarshal(event.Kv.Value, &podInfo)
		//utils.HandleError("unmarshal pod error", err)
		fmt.Println(podInfo)
		output, _ := json.Marshal(podInfo)
		fmt.Printf("%v\n", string(output))
		CreateAndRunPod(podInfo)
	}
	return err
}
