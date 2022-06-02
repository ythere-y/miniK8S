package controllerManager

import (
	"encoding/json"
	"fmt"
	"minik8s/apiserver"
	pod2 "minik8s/registry/pod"

	"minik8s/utils"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	yaml "gopkg.in/yaml.v2"
)

func createPod(event *clientv3.Event) error {
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
		err = json.Unmarshal(event.Kv.Value, &keys)
		if err != nil {
			return err
		}

		apiserver.DeletePods(keys)
	}
	return err
}

func stopPod(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		var keys []string
		err = json.Unmarshal(event.Kv.Value, &keys)
		if err != nil {
			return err
		}
		apiserver.StopPods(keys)
	}
	return err
}
