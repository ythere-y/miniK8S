package controllerManager

import (
	"fmt"
	"minik8s/apiserver"
	"minik8s/pod"
	"minik8s/replicaset"

	"minik8s/apiserver"
	"minik8s/pod"
	"minik8s/utils"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	yaml "gopkg.in/yaml.v2"
)

func podControllerManagerHandler(event *clientv3.Event) error {
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

func rsControllerManagerHandler(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		var newRsYaml replicaset.RSyaml
		err = yaml.Unmarshal(event.Kv.Value, &newRsYaml)
		if err != nil {
			fmt.Printf("yaml unmarshal error->:\n%v\n", err.Error())
		}
		rsInfo := replicaset.YamlToRS(newRsYaml)
		err = apiserver.SaveRsInfo(rsInfo)
		utils.HandleError("save replicaset info error", err)
		// TODO: deal with pod replicas, create pod in some nodes

	}
}
