package controllerManager

import (
	"encoding/json"
	"fmt"
	"minik8s/apiserver"
	"minik8s/constant"
	"minik8s/lab/etcd"
	. "minik8s/lab/etcd"
	pod2 "minik8s/registry/pod"
	"minik8s/utils"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v2"
)

func podControllerWatch() {
	fmt.Println("[Pod controller] init!")
	var watchName string
	watchName = SetKey(
		SetPrefix(constant.ControllerPrefix),
		SetSourceType(constant.PodSourceName),
		JustAppend(constant.CREATE),
	)
	apiserver.SyncWatch(watchName, createPod)

	watchName = SetKey(
		SetPrefix(constant.ControllerPrefix),
		SetSourceType(constant.PodSourceName),
		JustAppend(constant.DELETE),
	)
	apiserver.SyncWatch(watchName, deletePod)

	watchName = SetKey(
		SetPrefix(constant.ControllerPrefix),
		SetSourceType(constant.PodSourceName),
		JustAppend(constant.STOP),
	)
	apiserver.SyncWatch(watchName, stopPod)

	watchName = SetKey(
		SetPrefix(constant.ControllerPrefix),
		SetSourceType(constant.PodSourceName),
		JustAppend(constant.WATCH),
	)
	apiserver.SyncWatch(watchName, watchPod)

	watchName = SetKey(
		SetPrefix(constant.WatchPrefix),
		SetSourceType(constant.PodSourceName),
	)
	apiserver.SyncWatch(watchName, dealFail)
}

// region 增

func createPod(event *clientv3.Event) error {
	fmt.Printf("controller start creating pod\n")
	var err error
	switch event.Type {
	case mvccpb.PUT:
		var newPodYaml pod2.PodYaml
		err = yaml.Unmarshal(event.Kv.Value, &newPodYaml)
		if err != nil {
			fmt.Printf("yaml unmarshal error->:\n%v\n", err.Error())
		}
		podInfo := pod2.YamlToPod(newPodYaml)
		podName := podInfo.Meta.Name

		for _, pod := range MemPods {
			if pod.Meta.Name == podName {
				fmt.Printf("pod %v already exist~!\n", podName)
				return err
			}
		}

		AddPod(podInfo)
		err = apiserver.SavePodInfo(podInfo)
		utils.HandleError("save pod info error", err)

		// 这里取消了scheduler的数据转发层，直接把scheduler作为controller内部的一个分支
		err = DisPodtoNode(podName)
		// err = apiserver.PushPodToScheduler(podInfo)
		// utils.HandleError("push pod to scheduler error", err)
	}
	return err
}

// endregion

// region 删

func stopPod(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		var keys []string
		err = json.Unmarshal(event.Kv.Value, &keys)
		if err != nil {
			return err
		}
		apiserver.ActStopPods(keys)
	}
	return err
}

// TODO  delete的步骤还需要先stop，还需要考虑如何stop结束之后再delete
func deletePod(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		var keys []string
		err = json.Unmarshal(event.Kv.Value, &keys)
		if err != nil {
			return err
		}
		RemovePods(keys)
		apiserver.ActDeletePods(keys)
	}
	return err

}

func watchPod(event *clientv3.Event) error {
	fmt.Printf("controller start watching pod\n")
	// wait for pod creating
	time.Sleep(12 * time.Second)
	fmt.Println("controller watch pod finish sleep")
	var err error
	switch event.Type {
	case mvccpb.PUT:
		var podName string
		podName = string(event.Kv.Value)
		// 根据Podname获得node name
		nodeName, err := etcd.GetValueWithPrefix(
			etcd.SetKey(
				etcd.SetPrefix(constant.RelationPrefix),
				etcd.SetSourceType(constant.PodSourceName),
				etcd.JustAppend(podName),
				etcd.JustAppend(constant.NodeSourceName)))
		utils.HandleError("watchPod controller get node value error", err)
		// 通知kubelet监控健康状态
		buildKey := etcd.SetKey(
			etcd.SetPrefix(constant.WatchPrefix),
			etcd.SetSourceType(constant.NodeSourceName),
			etcd.JustAppend(nodeName),
			etcd.JustAppend(podName))
		buildvalue := podName
		apiserver.SyncPut(buildKey, buildvalue)
	}
	return err
}

// node节点上有Pod fail了
func dealFail(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		if RsRunningFlag {
			podName := string(event.Kv.Value)
			// //先删除旧的Pod
			// apiserver.CmdDeletePod([]string{podName})
			//重新创建fail的pod, 就分配到原node上
			nodeName, err := etcd.GetValueWithPrefix(
				etcd.SetKey(
					etcd.SetPrefix(constant.RelationPrefix),
					etcd.SetSourceType(constant.PodSourceName),
					etcd.JustAppend(podName),
					etcd.JustAppend(constant.NodeSourceName)))
			utils.HandleError("dealFail controller get node value error", err)
			err = apiserver.DistributePodtoNode(nodeName, podName)
			utils.HandleError("Distribute pod to node error", err)
		}
		err = nil
	}
	return err
}

//CmdStopPods

// endregion

// region 改

// endregion

// region 查

// endregion

// region 送

// endregion
