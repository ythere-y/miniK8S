package apiserver

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"io/ioutil"
	"minik8s/constant"
	"minik8s/lab/etcd"
	pod2 "minik8s/registry/pod"
	"minik8s/utils"
	"time"
)

// region 增

func CmdCreatePod(filename string) {
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
func SavePodInfo(pod pod2.Pod) error {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType(constant.PodSourceName),
		etcd.SetPodName(pod.Meta.Name))
	value, err := json.Marshal(pod)

	utils.HandleError("marshal pod error", err)
	etcd.Put(key, string(value))
	return err
}

// endregion

// region 删

func ActStopPods(podsNames []string) {

	var deleteTargets []string
	var values []string

	for _, key := range podsNames {
		buildKey := etcd.SetKey(
			etcd.SetPrefix(constant.RelationPrefix),
			etcd.SetSourceType(constant.PodSourceName),
			etcd.SetPodName(key))
		nodeName, err := etcd.GetValue(buildKey)
		utils.HandleError("get value error", err)

		buildKey = etcd.SetKey(
			etcd.SetPrefix(constant.RelationPrefix),
			etcd.SetSourceType(constant.NodeSourceName),
			etcd.SetNodeName(nodeName),
			etcd.SetPodName(key),
		)
		deleteTargets = append(deleteTargets, buildKey)
		values = append(values, constant.StopFlag)
	}
	for _, del := range deleteTargets {
		fmt.Printf("stop [key = %v]\n", del)
	}
	SyncPutList(deleteTargets, values)
}

//CmdStopPods
/*
按照podName让pod停止运行（container的状态变为stopped）
主要流程是先找到，然后stop，然后delete
*/
func CmdStopPods(names []string) {
	nameSet := utils.ParseNames(names)
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

//CmdDeletePod
/*
按照podName删除一个pod
主要流程是先找到，然后stop，然后delete
*/
func CmdDeletePod(names []string) {
	// 检查是否存在在relation关系中
	nameSet := utils.ParseNames(names)
	for _, name := range nameSet {
		path := etcd.SetKey(
			etcd.SetPrefix(constant.RelationPrefix),
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

func ActDeletePods(podsNames []string) {
	var (
		deleteTargets []string
		buildKey      string
		err           error
		nodeName      string
		serviceName   string
	)
	for _, key := range podsNames {

		// 先查询relations找到所属的node
		buildKey = etcd.SetKey(
			etcd.SetPrefix(constant.RelationPrefix),
			etcd.SetSourceType(constant.PodSourceName),
			etcd.SetPodName(key),
			etcd.SetSourceType(constant.NodeSourceName))
		nodeName, err = etcd.GetValue(buildKey)
		utils.HandleError("get value error", err)
		deleteTargets = append(deleteTargets, buildKey)
		if nodeName != "" {
			buildKey = etcd.SetKey(
				etcd.SetPrefix(constant.RelationPrefix),
				etcd.SetSourceType(constant.NodeSourceName),
				etcd.SetNodeName(nodeName),
				etcd.SetPodName(key),
			)
			deleteTargets = append(deleteTargets, buildKey)
		}

		// 先查询relations找到对应的service
		buildKey = etcd.SetKey(
			etcd.SetPrefix(constant.RelationPrefix),
			etcd.SetSourceType(constant.PodSourceName),
			etcd.SetPodName(key),
			etcd.SetSourceType(constant.ServiceSourceName))
		serviceName, err = etcd.GetValue(buildKey)
		utils.HandleError("get value error", err)
		deleteTargets = append(deleteTargets, buildKey)
		if serviceName != "" {
			buildKey = etcd.SetKey(
				etcd.SetPrefix(constant.RelationPrefix),
				etcd.SetSourceType(constant.ServiceSourceName),
				etcd.SetNodeName(serviceName),
				etcd.SetPodName(key),
			)
			deleteTargets = append(deleteTargets, buildKey)
		}
		buildKey = etcd.SetKey(
			etcd.SetPrefix(constant.RelationPrefix),
			etcd.SetSourceType(constant.PodSourceName),
			etcd.SetPodName(key))
		deleteTargets = append(deleteTargets, buildKey)

		buildKey = etcd.SetKey(
			etcd.SetPrefix(constant.RegistryPrefix),
			etcd.SetSourceType(constant.PodSourceName),
			etcd.SetPodName(key))
		deleteTargets = append(deleteTargets, buildKey)

	}
	for _, del := range deleteTargets {
		fmt.Printf("del [key = %v]\n", del)
	}
	SyncDel(deleteTargets)
}

// endregion

// region 改

func UpdatePodToKubelet(podName string) error {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RelationPrefix),
		etcd.SetSourceType(constant.NodeSourceName),
		etcd.SetPodName(podName),
	)
	value := constant.UpdateFlag
	SyncPut(key, value)
	return nil
}

//SetPodsStatus
/*
修改Pods状态
*/
func SetPodsStatus(podName string, targetPod pod2.Pod) {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType(constant.PodSourceName),
		etcd.SetPodName(podName))
	value, _ := json.Marshal(targetPod)
	etcd.Put(key, string(value))
}

// endregion

// region 查

//GetPodInfo
/*
 根据podName，从/registry/pods/目录下寻找对应的pod
并组建成Pod返回
*/
func GetPodInfo(podName string) *pod2.Pod {
	getRes, err := etcd.Get(etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType("pods"),
		etcd.SetName(podName)))

	utils.HandleError("get pod info error", err)
	if len(getRes) == 0 {
		return nil
	}
	var podInfo pod2.Pod
	err = json.Unmarshal([]byte(getRes[0]), &podInfo)
	utils.HandleError("unmarshal pod failed ", err)

	return &podInfo
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
	pod2.PodPreDisplay()
	for _, event := range getRsp.Kvs {
		tmpValue := event.Value
		var tmpPod pod2.Pod
		err = json.Unmarshal(tmpValue, &tmpPod)
		utils.HandleError("unmarshal pod error", err)
		tmpPod.Display()
	}
}

//DisplayPodInfo
/*
按照podName查找Pods的信息（以表格形式打印主要信息）
*/
func DisplayPodInfo(name string) {
	tmpPod := GetPodInfo(name)
	if tmpPod == nil {
		fmt.Println("cannot find any pods")
	} else {
		pod2.PodPreDisplay()
		tmpPod.Display()
	}
}

// endregion

// region 送

func PushPodToScheduler(pod pod2.Pod) error {

	key := etcd.SetKey(
		etcd.SetPrefix(constant.SchedulerPrefix),
		etcd.SetSourceType(constant.PodSourceName),
		etcd.SetName(pod.Meta.Name))

	value := pod.Meta.Name
	SyncPut(key, value)

	return nil
}

// endregion
