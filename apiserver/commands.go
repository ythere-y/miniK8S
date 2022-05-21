package apiserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"minik8s/constant"
	"minik8s/lab/etcd"
	"minik8s/pod"
	"minik8s/replicaset"
	"minik8s/utils"

	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func CreatePod(filename string) {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.ControllerPrefix),
		etcd.JustAppend(constant.PodSourceName),
		etcd.JustAppend(constant.CREATE),
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
		etcd.SetSourceType(constant.PodSourceName),
		etcd.SetPodName(pod.Meta.Name))
	value, err := json.Marshal(pod)

	utils.HandleError("marshal pod error", err)
	etcd.Put(key, string(value))
	return err
}

func PushPodToScheduler(pod pod.Pod) error {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.SchedulerPrefix),
		etcd.SetSourceType(constant.PodSourceName),
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
		etcd.SetSourceType(constant.NodeSourceName),
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
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
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

//parseNames
/*
解析多个名字，去除其中的重复内容
*/
func parseNames(names []string) []string {
	var nameSet []string
	for index, name := range names {
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
	}
	return nameSet
}

//DeletePod
/*
按照podName删除一个pod
主要流程是先找到，然后stop，然后delete
*/
func DeletePod(names []string) {
	nameSet := parseNames(names)
	for _, name := range nameSet {
		path := etcd.SetKey(
			etcd.SetPrefix(constant.RegistryPrefix),
			etcd.SetSourceType(constant.PodSourceName),
			etcd.SetPodName(name))
		if CheckIfExist(path) == false {
			fmt.Println("cannot find pod [" + name + "] !")
			return
		}
	}
	names = nameSet

	key := etcd.SetKey(
		etcd.SetPrefix(constant.ControllerPrefix),
		etcd.SetSourceType(constant.PodSourceName),
		etcd.JustAppend(constant.DELETE),
		etcd.JustAppend(time.Now().String()))
	value, err := json.Marshal(names)
	if err != nil {
		panic(err)
		return
	}
	SyncPut(key, string(value))
}

//StopPod
/*
按照podName让pod停止运行（container的状态变为stopped）
主要流程是先找到，然后stop，然后delete
*/
func StopPod(names []string) {
	nameSet := parseNames(names)
	for _, name := range nameSet {
		if CheckPodIfExist(name) == false {
			fmt.Println("cannot find pod [" + name + "] !")
			return
		}
	}
	names = nameSet

	key := etcd.SetKey(
		etcd.SetPrefix(constant.ControllerPrefix),
		etcd.SetSourceType(constant.PodSourceName),
		etcd.JustAppend(constant.STOP),
		etcd.JustAppend(time.Now().String()))
	value, err := json.Marshal(names)
	if err != nil {
		panic(err)
		return
	}
	SyncPut(key, string(value))
}

// Replicaset related functions
/*
 * create a replicaset
 * input: yaml file
 */
func CreateRs(file string) {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.ControllerPrefix),
		etcd.JustAppend(constant.ReplicaSourceName),
		etcd.JustAppend(constant.CREATE),
		etcd.JustAppend(time.Now().String()))
	value, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("file %v read error!\n", file)
	}
	SyncPut(key, string(value))
}

func SaveRsInfo(rs replicaset.ReplicaSet) error {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType(constant.ReplicaSourceName),
		etcd.SetPodName(rs.RSmeta.Name))
	value, err := json.Marshal(rs)

	utils.HandleError("marshal pod error", err)
	etcd.Put(key, string(value))
	return err
}

/*
 * pods managed by replicaset should be put into a new key
 * in etcd, which is /registry/rspods/[rs_name]/[pod_name]
 */
func SaveRsPodInfo(rs replicaset.ReplicaSet, pod pod.Pod) error {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType(constant.RsPodSourceName),
		etcd.JustAppend(rs.RSmeta.Name),
		etcd.JustAppend(pod.Meta.Name),
	)
	value, err := json.Marshal(pod)
	utils.HandleError("marshal rspod error", err)
	etcd.Put(key, string(value))
	return err
}
