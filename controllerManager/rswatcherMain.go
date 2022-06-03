package controllerManager

import (
	"encoding/json"
	"minik8s/apiserver"
	"minik8s/constant"
	"minik8s/lab/etcd"
	"minik8s/registry/replicaset"
	"minik8s/utils"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// watch /registry/replicaset
func RsWatch() {
	// var watchName string
	watchName := etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType(constant.ReplicaSourceName),
	)
	apiserver.SyncWatch(watchName, watchRs)
}

// watch /relations/replicaset/[rsname]
func RsPodWatch(rsName string) {
	watchName := etcd.SetKey(
		etcd.SetPrefix(constant.RelationPrefix),
		etcd.SetSourceType(constant.ReplicaSourceName),
		etcd.JustAppend(rsName))
	apiserver.SyncWatch(watchName, watchRsPod)
}

func ListenRpod(podName string) {
	watchName := etcd.SetKey(
		etcd.SetPrefix(constant.RelationPrefix),
		etcd.SetSourceType(constant.ReplicaSourceName),
		etcd.JustAppend(podName))
	apiserver.SyncWatch(watchName, watchRsPod)
}

// 监听rs创立，每新建一个rs，就对其pod进行监听
func watchRs(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		var rsInfo replicaset.ReplicaSet
		err := json.Unmarshal(event.Kv.Value, &rsInfo)
		utils.HandleError("watch rs json unmarshal error", err)
		rsName := rsInfo.RSmeta.Name
		// 对对应的pod位置进行监听
		RsPodWatch(rsName)
	}
	return err
}

//对replicaset里面的pod进行监听处理
//交给pod controller处理，value为需要监听的pod name
func watchRsPod(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		var podName string
		podName = string(event.Kv.Value)
		watchKey := etcd.SetKey(
			etcd.SetPrefix(constant.ControllerPrefix),
			etcd.SetSourceType(constant.PodSourceName),
			etcd.JustAppend(constant.WATCH),
			etcd.JustAppend(time.Now().String()))
		value := podName
		apiserver.SyncPut(watchKey, string(value))
	}
	return err
}
