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
	"minik8s/utils"
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

	case mvccpb.PUT:
		// 增加一个services的操作
		fmt.Printf("\n\n[**]kubelet handling key = %v, value = %v\n", string(event.Kv.Key), string(event.Kv.Value))
		var servieInfo service.Service
		operation := utils.GetLastWord(string(event.Kv.Key))
		fmt.Printf("[***] the operation is %v\n", operation)
		err = json.Unmarshal(event.Kv.Value, &servieInfo)

		switch operation {
		case constant.CREATE:
			fmt.Printf("[set]\n")
			// 遍历写入iptables
			for _, podip := range servieInfo.PodIp {
				err = environment.SetIptables(servieInfo.ServiceIP, servieInfo.ServicePort, podip, servieInfo.TargetPort)
				if err != nil {
					panic(err)
				}
			}
		case constant.DELETE:
			fmt.Printf("[remove]\n")
			// 遍历写入iptables
			for _, podip := range servieInfo.PodIp {
				err = environment.RemoveIptables(servieInfo.ServiceIP, servieInfo.ServicePort, podip, servieInfo.TargetPort)
				if err != nil {
					panic(err)
				}
			}

		}

	}
	return err
}
