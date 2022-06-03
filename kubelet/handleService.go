package kubelet

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"minik8s/apiserver"
	"minik8s/constant"
	"minik8s/environment"
	. "minik8s/lab/etcd"
	"minik8s/registry/service"
)

func kubeleteServiceWatch() {

	watchName := SetKey(
		SetPrefix(constant.KubeletPrefix),
		SetSourceType(constant.ServiceSourceName),
	)
	apiserver.SyncWatch(watchName, servicehandler)

}
func servicehandler(event *clientv3.Event) error {
	var err error

	//serviceName := utils.GetLastWord(string(event.Kv.Key))

	switch event.Type {

	case mvccpb.DELETE:
		// 删除一个service的操作
		fmt.Printf("kubelet handling Delete service = %v\n", string(event.Kv.Key))
		var servieInfo service.Service
		// 遍历写入iptables
		for _, podip := range servieInfo.PodIp {
			err = environment.RemoveIptables(servieInfo.ServiceIP, servieInfo.ServicePort, podip, servieInfo.TargetPort)
			if err != nil {
				panic(err)
			}
		}

	case mvccpb.PUT:
		// 增加一个services的操作
		fmt.Printf("kubelet handling key = %v, value = %v\n", string(event.Kv.Key), string(event.Kv.Value))
		var servieInfo service.Service

		// 遍历写入iptables
		for _, podip := range servieInfo.PodIp {
			err = environment.SetIptables(servieInfo.ServiceIP, servieInfo.ServicePort, podip, servieInfo.TargetPort)
			if err != nil {
				panic(err)
			}
		}

		err = json.Unmarshal(event.Kv.Value, &servieInfo)

	}
	return err
}
