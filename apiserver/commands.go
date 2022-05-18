package apiserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"minik8s/constant"
	"minik8s/lab/etcd"
	"minik8s/src/pod"
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

func PushPodToScheduler(pod pod.Pod) error {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.SchedulerPrefix),
		etcd.JustAppend("pods"))
	value, err := json.Marshal(pod)
	if err != nil {
		fmt.Printf("Marshal pod into json error -> :\n %v\n", err.Error())
		return err
	}
	etcd.SyncPut(key, string(value))
	return err
}

func DistributePodtoNode(nodeName string, pod pod.Pod) error {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType("pods"),
		etcd.SetNameSpace("default"),
		etcd.SetNodeName(nodeName),
		etcd.SetPodName("pod_1"))

	value, err := json.Marshal(pod)
	if err != nil {
		fmt.Printf("Marshal pod error ->:\n%v\n", err.Error())
	}
	etcd.SyncPut(key, string(value))

	return err
}
