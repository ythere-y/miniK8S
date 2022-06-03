package apiserver

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"minik8s/constant"
	"minik8s/lab/etcd"
	"minik8s/utils"
	"time"
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

//SyncWatchWithTime
/*
启动watch name，使用前缀watch
有变动之后使用handler函数处理，超时之后不再watch
*/
func SyncWatchWithTime(name string, handler etcd.Handler, fail etcd.FailOut, timeset time.Duration) {
	go etcd.WatchWithFuncWithTime(name, handler, fail, timeset)
}

//SerlWatchWithTime
/*
启动watch name，使用前缀watch
串行watch，会阻塞当前线程
有变动之后使用handler函数处理，并且有超时设置，超时之后不再阻塞
*/
func SerlWatchWithTime(name string, handler etcd.Handler, fail etcd.FailOut, timeset time.Duration) {
	go etcd.WatchWithFuncWithTime(name, handler, fail, timeset)

}

//SyncWatch
/*
启动watch name，使用前缀watch
有变动之后使用handler函数处理
*/
func SyncWatch(name string, handler etcd.Handler) {
	go etcd.WatchWithFunc(name, handler)
}

func SerlPut(key string, value string) {
	etcd.Put(key, value)
}

//SyncPut
/*
向etcd中put一个k-v对
*/
func SyncPut(key string, value string) {
	go etcd.Put(key, value)
}

func Reply(key []byte, reply string) {
	SyncPut(
		etcd.SetKey(
			etcd.SetPrefix(constant.ReplayPrefix),
			etcd.JustAppend(utils.GetLastWord(string(key)))),
		reply)
}

//SyncPutList
/*
向etcd中put一连串的k-v对
*/
func SyncPutList(key []string, value []string) {
	go etcd.PutList(key, value)
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

// region 添加关系relation

func SaveRelationPodtoNode(keyname string, valname string) {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RelationPrefix),
		etcd.SetSourceType(constant.PodSourceName),
		etcd.SetPodName(keyname),
		etcd.SetNodeName(valname))
	val := valname
	etcd.Put(key, val)
}

func SaveRelationNodetoPod(keyname string, valname string) {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RelationPrefix),
		etcd.SetSourceType(constant.NodeSourceName),
		etcd.SetPodName(keyname),
		etcd.SetNodeName(valname))
	val := valname
	etcd.Put(key, val)
}
func SaveRelationServicetoPod(keyname string, valname string) {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RelationPrefix),
		etcd.SetSourceType(constant.ServiceSourceName),
		etcd.SetPodName(keyname),
		etcd.SetNodeName(valname))
	val := valname
	etcd.Put(key, val)
}

// endregion

// region 删除关系

func DeleteRelationPodtoNode(keyname string, valname string) {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RelationPrefix),
		etcd.SetSourceType(constant.PodSourceName),
		etcd.SetPodName(keyname),
		etcd.SetNodeName(valname))
	etcd.Delete(key)
}
func DeleteRelationNodetoPod(keyname string, valname string) {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RelationPrefix),
		etcd.SetSourceType(constant.ServiceSourceName),
		etcd.SetPodName(keyname),
		etcd.SetNodeName(valname))
	etcd.Delete(key)
}

func DeleteRelationServicetoPod(keyname string, valname string) {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RelationPrefix),
		etcd.SetSourceType(constant.ServiceSourceName),
		etcd.SetPodName(keyname),
		etcd.SetNodeName(valname))
	etcd.Delete(key)
}

// endregion
