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
		etcd.JustAppend("pods"),
		etcd.JustAppend("create"),
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

//GetPodInfo
/*
 根据podName，从/registry/pods/目录下寻找对应的pod
并组建成Pod返回
*/
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

//GetChildNum
/*
根据key值，使用前缀查找，统计查询到的结果数量并返回
*/
func GetChildNum(key string) uint32 {
	var (
		num    uint32 = 0
		getRsp *clientv3.GetResponse
		err    error
	)
	num = 1
	getRsp, err = etcd.GetWithPrefix(key)
	utils.HandleError("get with prefix error[from get child num]", err)
	num = uint32(len(getRsp.Kvs))

	return num
}

//CheckIfExist
/*
根据key值，使用准确查找，检查某个key值是否存在
*/
func CheckIfExist(key string) bool {
	var (
		exist  bool = false
		getRsp *clientv3.GetResponse
		err    error
	)
	getRsp, err = etcd.GetNormal(key)
	utils.HandleError("get with prefix error[from get child num]", err)
	exist = len(getRsp.Kvs) == 1
	return exist

}

//SyncWatch
/*
启动watch name，使用前缀watch
有变动之后使用handler函数处理
*/
func SyncWatch(name string, handler etcd.Handler) {
	go etcd.WatchWithFunc(name, handler)
}

//SyncPut
/*
向etcd中put一个k-v对
*/
func SyncPut(key string, value string) {
	go etcd.Put(key, value)
}

//SyncPutList
/*
向etcd中put一连串的k-v对
*/
func SyncPutList(key []string, value []string) {
	go etcd.PutList(key, value)
}
