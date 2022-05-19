package apiserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"minik8s/constant"
	"minik8s/lab/etcd"
	"minik8s/pod"
	"minik8s/utils"
	"strconv"
)

var command_id = 1

func CreatePod(filename string) {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.ControllerPrefix),
		etcd.JustAppend("pods"),
		etcd.JustAppend("id_"+strconv.Itoa(command_id)))
	value, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("file %v read error!\n", filename)
	}
	etcd.SyncPut(key, string(value))
}

func SavePodInfo(pod pod.Pod) error {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType("pods"),
		etcd.SetPodName(pod.Meta.Name))
	value, err := json.Marshal(pod)

	utils.HandleError("marshal pod error", err)
	etcd.Put(key, string(value))
	return err
}
func PushPodToScheduler(pod pod.Pod) error {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.SchedulerPrefix),
		etcd.JustAppend("pods"),
		etcd.SetName(pod.Meta.Name))

	value := pod.Meta.Name
	etcd.SyncPut(key, value)
	return nil
}

func DistributePodtoNode(nodeName string, podName string) error {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType("pods"),
		etcd.SetNameSpace("default"),
		etcd.SetNodeName(nodeName),
		etcd.SetPodName(podName))

	value := podName
	etcd.SyncPut(key, value)

	return nil
}

func GetPodInfo(podName string) pod.Pod {
	getRes, err := etcd.Get(etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType("pods"),
		etcd.SetName(podName)))

	utils.HandleError("get pod info error", err)

	var podInfo pod.Pod
	err = json.Unmarshal([]byte(getRes[0]), &podInfo)
	utils.HandleError("unmarshal pod failed ", err)

	return podInfo
}
