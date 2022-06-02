package controllerManager

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v2"
	"minik8s/apiserver"
	"minik8s/constant"
	. "minik8s/lab/etcd"
	pod2 "minik8s/registry/pod"
	"minik8s/utils"
)

func init() {
	fmt.Printf("Pod controller init!")
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
}

// region 增

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

		for _, pod := range MemPods {
			if pod.Meta.Name == podName {
				fmt.Printf("pod %v already exist~!\n", podName)
				return err
			}
		}

		AddPod(podInfo)
		err = apiserver.SavePodInfo(podInfo)
		utils.HandleError("save pod info error", err)
		err = apiserver.PushPodToScheduler(podInfo)
		utils.HandleError("push pod to scheduler error", err)
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

//CmdStopPods

// endregion

// region 改

// endregion

// region 查

// endregion

// region 送

// endregion
