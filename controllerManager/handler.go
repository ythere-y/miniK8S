package controllerManager

import (
	"encoding/json"
	"fmt"
	"minik8s/apiserver"
	"minik8s/constant"
	"minik8s/lab/etcd"
	"minik8s/pod"
	"minik8s/replicaset"
	"reflect"
	"strconv"

	"minik8s/apiserver"
	"minik8s/pod"
	"minik8s/utils"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	yaml "gopkg.in/yaml.v2"
)

/*
 * Pod controller handler, not used
 */
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

/*
 * Replicaset controller handler, not used
 */
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
		// TODO: deal with pod replicas, create pod in some nodes. UNFINISHED
		if reflect.DeepEqual(rsInfo.RSspec.SelectorLabels,
			rsInfo.PodTemplate.Meta.Labels) {
			replicas := rsInfo.RSspec.Replicas
			pods := replicaset.CreatePodInstances(rsInfo, replicas)
			for _, pod := range pods {
				apiserver.SaveRsPodInfo(rsInfo, pod)
				apiserver.PushPodToScheduler(pod)
			}
		}
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
		err = apiserver.SavePodInfo(podInfo)
		utils.HandleError("save pod info error", err)
		err = apiserver.PushPodToScheduler(podInfo)
		utils.HandleError("push pod to scheduler error", err)
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

func createReplicaset(event *clientv3.Event) error {
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
		// TODO: deal with pod replicas, create pod in some nodes.
		if reflect.DeepEqual(rsInfo.RSspec.SelectorLabels,
			rsInfo.PodTemplate.Meta.Labels) {
			replicas := rsInfo.RSspec.Replicas
			pods := replicaset.CreatePodInstances(rsInfo, replicas)
			for _, pod := range pods {
				// apiserver.SaveRsPodInfo(rsInfo, pod)
				err = apiserver.SavePodInfo(pod)
				utils.HandleError("rs save pod info error", err)
				err = apiserver.PushPodToScheduler(pod)
				utils.HandleError("rs push pod to scheduler error", err)
			}
		}
	}
	return err
}

func deleteReplicaset(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		var names []string
		var rsDelTar []string
		var podDelTar []string
		err = json.Unmarshal(event.Kv.Value, &names)
		for _, rsname := range names {
			// set rs key to be deleted
			rskey := etcd.SetKey(
				etcd.SetPrefix(constant.RegistryPrefix),
				etcd.SetSourceType(constant.ReplicaSourceName),
				etcd.JustAppend(rsname))
			rsDelTar = append(rsDelTar, rskey)

			// set pod key to be deleted
			rsInfo := apiserver.GetRsInfo(rsname)
			replicas := rsInfo.RSspec.Replicas
			podname := rsInfo.PodTemplate.Meta.Name
			for i := 1; i <= replicas; i++ {
				rspodname := podname + "-" + strconv.Itoa(i)
				// 之后会交给deletePod处理
				podDelTar = append(podDelTar, rspodname)
			}
		}
		// delete rs info
		apiserver.SyncDel(rsDelTar)
		// delete pods in these rs
		apiserver.DeletePod(podDelTar)
	}
	return err
}
