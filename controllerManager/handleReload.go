package controllerManager

import (
	"encoding/json"
	"minik8s/apiserver"
	"minik8s/constant"
	"minik8s/lab/etcd"
	"minik8s/registry/pod"
	"minik8s/utils"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func ReloadControllerWatch() {
	var watchName string
	watchName = etcd.SetKey(
		etcd.SetPrefix(constant.ControllerPrefix),
		etcd.SetSourceType(constant.RELOAD),
	)
	apiserver.SyncWatch(watchName, sysReload)
}

func sysReload(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
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
			return err
		}
		for _, event := range getRsp.Kvs {
			tmpValue := event.Value
			var tmpPod pod.Pod
			err = json.Unmarshal(tmpValue, &tmpPod)
			podName := tmpPod.Meta.Name
			utils.HandleError("unmarshal pod error", err)
			// get its node
			nodeName, err := etcd.GetValueWithPrefix(
				etcd.SetKey(
					etcd.SetPrefix(constant.RelationPrefix),
					etcd.SetSourceType(constant.PodSourceName),
					etcd.JustAppend(podName),
					etcd.JustAppend(constant.NodeSourceName)))
			utils.HandleError("dealFail controller get node value error", err)
			// re run
			err = apiserver.DistributePodtoNode(nodeName, podName)
			utils.HandleError("Distribute pod to node error", err)
		}
	}
	return err
}
