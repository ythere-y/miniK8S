package apiserver

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"io/ioutil"
	"minik8s/constant"
	"minik8s/lab/etcd"
	pod2 "minik8s/registry/pod"
	"minik8s/registry/service"
	"minik8s/utils"
	"time"
)

// region 增

func SaveServiceInfo(service service.Service) error {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType(constant.PodSourceName),
		etcd.SetPodName(service.Name))
	value, err := json.Marshal(service)

	utils.HandleError("marshal service error", err)
	etcd.Put(key, string(value))
	return err
}

func CreateServiceByFile(filename string) {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.ControllerPrefix),
		etcd.JustAppend(constant.ServiceSourceName),
		etcd.JustAppend(constant.CREATE),
		etcd.JustAppend(time.Now().String()))
	value, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Printf("file %v read error!\n", filename)
	}
	SyncPut(key, string(value))
}

func CreateService(status service.Service) {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType(constant.ServiceSourceName),
		etcd.SetNodeName(status.Name))
	value, err := json.Marshal(status)
	utils.HandleError("marshal service error", err)
	SyncPut(key, string(value))
}

// endregion

// region 删

func ActDeleteService(names []string) {
	// 1. 删除relations中的此service的目录
	// 2. 删除registry中次service的目录
	var deleteTargets []string
	for _, key := range names {
		// 先查询relations找到
		buildKey := etcd.SetKey(
			etcd.SetPrefix(constant.RelationPrefix),
			etcd.SetSourceType(constant.ServiceSourceName),
			etcd.SetPodName(key))
		deleteTargets = append(deleteTargets, buildKey)

		buildKey = etcd.SetKey(
			etcd.SetPrefix(constant.RegistryPrefix),
			etcd.SetSourceType(constant.ServiceSourceName),
			etcd.SetPodName(key))
		deleteTargets = append(deleteTargets, buildKey)

	}
	for _, del := range deleteTargets {
		fmt.Printf("del [key = %v]\n", del)
	}
	SyncDel(deleteTargets)
}

func CmdDeleteService(names []string) {
	// 检查是否存在在relation关系中
	nameSet := utils.ParseNames(names)
	for _, name := range nameSet {
		path := etcd.SetKey(
			etcd.SetPrefix(constant.RelationPrefix),
			etcd.SetSourceType(constant.ServiceSourceName),
			etcd.SetPodName(name))
		if CheckIfExist(path) == false {
			fmt.Println("cannot find service [" + name + "] !")
			return
		}
	}
	names = nameSet

	// 准备给controller发送命令
	key := etcd.SetKey(
		etcd.SetPrefix(constant.ControllerPrefix),
		etcd.SetSourceType(constant.ServiceSourceName),
		etcd.JustAppend(constant.DELETE),
		etcd.JustAppend(time.Now().String()))
	value, err := json.Marshal(names)
	if err != nil {
		panic(err)
		return
	}
	SyncPut(key, string(value))
}

// endregion

// region 改

func SetServiceStatus(name string, target service.Service) {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType(constant.ServiceSourceName),
		etcd.SetPodName(name))
	value, _ := json.Marshal(target)
	etcd.Put(key, string(value))
}

// endregion

// region 查

//GetNodeInfo
/*
 根据podName，从/registry/pods/目录下寻找对应的pod
并组建成Pod返回
*/
func GetServiceInfo(name string) *pod2.Pod {
	getRes, err := etcd.Get(etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType(constant.ServiceSourceName),
		etcd.SetName(name)))

	utils.HandleError("get service info error", err)
	if len(getRes) == 0 {
		return nil
	}
	var nodeInfo pod2.Pod
	err = json.Unmarshal([]byte(getRes[0]), &nodeInfo)
	utils.HandleError("unmarshal service failed ", err)

	return &nodeInfo
}

func DisplayAllServiceInfo() {
	var (
		num    uint32 = 0
		getRsp *clientv3.GetResponse
		err    error
	)
	getRsp, err = etcd.GetWithPrefix(
		etcd.SetKey(
			etcd.SetPrefix(constant.RegistryPrefix),
			etcd.SetSourceType(constant.ServiceSourceName)))
	utils.HandleError("get with prefix error[from get child num]", err)
	num = uint32(len(getRsp.Kvs))
	if num == 0 {
		fmt.Println("cannot find any service")
		return
	}
	service.ServicePreDisplay()
	for _, event := range getRsp.Kvs {
		tmpValue := event.Value
		var tmpPod service.Service
		err = json.Unmarshal(tmpValue, &tmpPod)
		utils.HandleError("unmarshal service error", err)
		tmpPod.Display()
	}
}

// DisplayServiceInfo
/*
按照podName查找Pods的信息（以表格形式打印主要信息）
*/
func DisplayServiceInfo(name string) {
	tmpService := GetServiceInfo(name)
	if tmpService == nil {
		fmt.Println("cannot find any node")
	} else {
		service.ServicePreDisplay()
		tmpService.Display()
	}
}

// endregion
