package controllerManager

import (
	"encoding/json"
	"fmt"
	"minik8s/apiserver"
	"minik8s/constant"
	"minik8s/lab/etcd"
	"minik8s/pod"

	"minik8s/utils"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	yaml "gopkg.in/yaml.v2"
)

func controllerManagerHandler(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		var newPodYaml pod.PodYaml
		err = yaml.Unmarshal(event.Kv.Value, &newPodYaml)
		if err != nil {
			fmt.Printf("yaml unmarshal error->:\n%v\n", err.Error())
		}
		podInfo := pod.YamlToPod(newPodYaml)
		err = apiserver.SavePodInfo(podInfo)
		utils.HandleError("save pod info error", err)
		err = apiserver.PushPodToScheduler(podInfo)
		utils.HandleError("push pod to scheduler error", err)
	}
	return err
}

func createPod(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		var newPodYaml pod.PodYaml
		err = yaml.Unmarshal(event.Kv.Value, &newPodYaml)
		if err != nil {
			fmt.Printf("yaml unmarshal error->:\n%v\n", err.Error())
		}
		podInfo := pod.YamlToPod(newPodYaml)
		podName := podInfo.Meta.Name
		err = apiserver.SavePodInfo(podInfo)
		utils.HandleError("save pod info error", err)
		if apiserver.CheckPodIfExist(podName) {
			// pod存在，换为更新操作
			err = apiserver.UpdatePodToKubelet(podName)
			utils.HandleError("update pod to kubelet error", err)
		} else {
			err = apiserver.PushPodToScheduler(podInfo)
			utils.HandleError("push pod to scheduler error", err)
		}
	}
	return err

}

func deletePod(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		var keys []string
		var deletargets []string
		err = json.Unmarshal(event.Kv.Value, &keys)
		if err != nil {
			return err
		}

		for _, key := range keys {
			buildKey := etcd.SetKey(
				etcd.SetPrefix(constant.RegistryPrefix),
				etcd.SetSourceType(constant.NodeSourceName),
				//TODO:需要一种能找到该pod归属于哪个node的机制
				etcd.SetNodeName("node_1"),
				etcd.SetPodName(key),
			)
			deletargets = append(deletargets, buildKey)

			buildKey = etcd.SetKey(
				etcd.SetPrefix(constant.RegistryPrefix),
				etcd.SetSourceType(constant.PodSourceName),
				etcd.SetPodName(key))
			deletargets = append(deletargets, buildKey)
		}
		for _, deletarget := range deletargets {
			fmt.Printf("del [key = %v]\n", deletarget)
		}
		apiserver.SyncDel(deletargets)
	}
	return err
}

func stopPod(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		var keys []string
		var stopTargets []string
		var values []string
		err = json.Unmarshal(event.Kv.Value, &keys)
		if err != nil {
			return err
		}

		for _, key := range keys {
			buildKey := etcd.SetKey(
				etcd.SetPrefix(constant.RegistryPrefix),
				etcd.SetSourceType(constant.NodeSourceName),
				//TODO:需要一种能找到该pod归属于哪个node的机制
				etcd.SetNodeName("node_1"),
				etcd.SetPodName(key),
			)
			stopTargets = append(stopTargets, buildKey)
			values = append(values, constant.StopFlag)
		}
		for _, deletarget := range stopTargets {
			fmt.Printf("stop [key = %v]\n", deletarget)
		}
		apiserver.SyncPutList(stopTargets, values)
	}
	return err
}
