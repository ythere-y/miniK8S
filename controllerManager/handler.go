package controllerManager

import (
	"fmt"
	"minik8s/apiserver"
	"minik8s/pod"

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
		err = apiserver.PushPodToScheduler(podInfo)
		if err != nil {
			fmt.Printf("Push Pod to Scheduler error -> :\n%v\n", err.Error())
		}
	}
	return err
}
