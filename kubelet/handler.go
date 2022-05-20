package kubelet

import (
	"fmt"
	"minik8s/constant"
	"minik8s/pod"
	"minik8s/utils"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"minik8s/apiserver"
)

func kubelethandler(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.DELETE:
		// 删除一个pod的操作
		fmt.Printf("kubelet handling Delete key = %v\n", string(event.Kv.Key))

		podName := utils.GetLastWord(string(event.Kv.Key))
		StopPod(podName)
		RemovePod(podName)
	case mvccpb.PUT:
		fmt.Printf("kubelet handling key = %v, value = %v\n", string(event.Kv.Key), string(event.Kv.Value))
		// 增加/修改 一个pod的操作
		podName := utils.GetLastWord(string(event.Kv.Key))
		operation := string(event.Kv.Value)
		switch operation {
		case podName:
			// 是创建操作
			var podInfo *pod.Pod
			podInfo = apiserver.GetPodInfo(podName)
			CreateAndRunPod(podInfo)
		case constant.StopFlag:
			// 是停止命令
			StopPod(podName)
		case constant.RemoveFlag:
			// 是删除命令
			StopPod(podName)
			RemovePod(podName)
		default:
			fmt.Printf("op = %v, it not in any!\n", operation)
		}
	}
	return err
}
