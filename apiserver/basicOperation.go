package apiserver

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"minik8s/constant"
	"minik8s/lab/etcd"
	"minik8s/utils"
)

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

func CheckPodIfExist(podName string) bool {
	var (
		exist  bool = false
		getRsp *clientv3.GetResponse
		err    error
		key    string
	)
	key = etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType(constant.PodSourceName),
		etcd.SetPodName(podName))
	getRsp, err = etcd.GetNormal(key)
	utils.HandleError("get with prefix error[from get child num]", err)
	exist = len(getRsp.Kvs) == 1
	return exist

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

func StopPods(podsNames []string) {

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

func DeletePods(podsNames []string) {
	var deleteTargets []string
	for _, key := range podsNames {

		// 先查询relations找到
		buildKey := etcd.SetKey(
			etcd.SetPrefix(constant.RelationPrefix),
			etcd.SetSourceType(constant.PodSourceName),
			etcd.SetPodName(key))
		nodeName, err := etcd.GetValue(buildKey)
		utils.HandleError("get value error", err)

		deleteTargets = append(deleteTargets, buildKey)
		buildKey = etcd.SetKey(
			etcd.SetPrefix(constant.RelationPrefix),
			etcd.SetSourceType(constant.NodeSourceName),
			//TODO:需要一种能找到该pod归属于哪个node的机制
			etcd.SetNodeName(nodeName),
			etcd.SetPodName(key),
		)
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

//SyncDel
/*
在etcd中删除某个key
*/
func SyncDel(key []string) {
	if len(key) == 0 {
		return
	} else if len(key) == 1 {
		go etcd.Delete(key[0])
	} else {
		go etcd.DeleteList(key)
	}
}
