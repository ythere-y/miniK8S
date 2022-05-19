package apiserver

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"io/ioutil"
	"minik8s/constant"
	"minik8s/lab/etcd"
	"minik8s/pod"
	"minik8s/utils"
	"time"
)

func CreatePod(filename string) {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.ControllerPrefix),
		etcd.JustAppend(constant.CREATE),
		etcd.JustAppend(constant.PodSourceName),
		etcd.JustAppend(time.Now().String()))
	value, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("file %v read error!\n", filename)
	}
	SyncPut(key, string(value))
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
	SyncPut(key, value)
	return nil
}

//DistributePodtoNode
/*
将一个pod分配给一个node（两者使用name来识别）
此为执行操作，将分配结果写入etcd
*/
func DistributePodtoNode(nodeName string, podName string) error {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType("nodes"),
		etcd.SetNameSpace("default"),
		etcd.SetNodeName(nodeName),
		etcd.SetPodName(podName))

	value := podName
	SyncPut(key, value)

	return nil
}

//SetPodsStatus
/*
修改Pods状态
*/
func SetPodsStatus(podName string, targetPod pod.Pod) {
	key := etcd.SetKey(etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType(constant.PodSourceName),
		etcd.SetPodName(podName))
	value, _ := json.Marshal(targetPod)
	etcd.Put(key, string(value))
}

//DisplayAllPodsInfo
/*
展示所有Pods的信息（以表格形式打印主要信息）
*/
func DisplayAllPodsInfo() {
	var (
		num    uint32 = 0
		getRsp *clientv3.GetResponse
		err    error
	)
	getRsp, err = etcd.GetWithPrefix(
		etcd.SetKey(
			etcd.SetPrefix(constant.RegistryPrefix),
			etcd.SetSourceType(constant.PodSourceName)))
	utils.HandleError("get with prefix error[from get child num]", err)
	num = uint32(len(getRsp.Kvs))
	if num == 0 {
		fmt.Println("cannot find any pods")
		return
	}
	pod.PodPreDisplay()
	for _, event := range getRsp.Kvs {
		tmpValue := event.Value
		var tmpPod pod.Pod
		err = json.Unmarshal(tmpValue, &tmpPod)
		utils.HandleError("unmarshal pod error", err)
		tmpPod.Display()
	}
}

//DisplayPodsInfo
/*
按照podName查找Pods的信息（以表格形式打印主要信息）
*/
func DisplayPodsInfo(name string) {
	tmpPod := GetPodInfo(name)
	if tmpPod == nil {
		fmt.Println("cannot find any pods")
	} else {
		pod.PodPreDisplay()
		tmpPod.Display()
	}
}

//DeletePod
/*
按照podName删除一个pod
主要流程是先找到，然后stop，然后delete
*/
func DeletePod(names []string) {
	var nameSet []string
	for index, name := range names {
		// 删除重名的
		haveTheSame := false
		for i := 0; i < index; i++ {
			if name == names[index] {
				haveTheSame = true
				break
			}
		}
		if haveTheSame == true {
			continue
		}
		nameSet = append(nameSet, name)
		targetPod := GetPodInfo(name)
		if targetPod == nil {
			fmt.Println("cannot find pod [" + name + "] !")
			return
		}

	}
	names = nameSet
	key := etcd.SetKey(
		etcd.SetPrefix(constant.ControllerPrefix),
		etcd.JustAppend(constant.DELETE),
		etcd.SetSourceType(constant.PodSourceName),
		etcd.JustAppend(time.Now().String()))
	value, err := json.Marshal(names)
	if err != nil {
		panic(err)
		return
	}
	SyncPut(key, string(value))
}
