package controllerManager

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v2"
	"minik8s/apiserver"
	"minik8s/constant"
	pod2 "minik8s/registry/pod"
	"minik8s/registry/service"

	. "minik8s/lab/etcd"
	"minik8s/utils"
)

func serviceControllerWatch() {
	fmt.Printf("Service controller init!")
	var watchName string
	watchName = SetKey(
		SetPrefix(constant.ControllerPrefix),
		SetSourceType(constant.ServiceSourceName),
		JustAppend(constant.CREATE),
	)
	apiserver.SyncWatch(watchName, createService)

	watchName = SetKey(
		SetPrefix(constant.ControllerPrefix),
		SetSourceType(constant.ServiceSourceName),
		JustAppend(constant.DELETE),
	)
	apiserver.SyncWatch(watchName, deleteService)

}

// region 增

func createService(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		var newone service.ServiceYaml
		err = yaml.Unmarshal(event.Kv.Value, &newone)
		if err != nil {
			fmt.Printf("yaml unmarshal error->:\n%v\n", err.Error())
		}
		serviceInfo := service.ServiceYamlToService(newone)
		//podInfo := pod2.YamlToPod(newone)
		serviceName := serviceInfo.Name

		// 检查重名
		for _, memService := range MemServices {
			if memService.Name == serviceName {
				fmt.Printf("service %v already exist !", serviceName)
				return err
			}
		}

		AddService(serviceInfo)
		err = apiserver.SaveServiceInfo(serviceInfo)
		utils.HandleError("save service info error", err)

		//1: label筛选
		var podSet []pod2.Pod
		if serviceInfo.Selector == nil {
			fmt.Printf("service selector should not be nil~!\n")
			return err
		}
		for _, pod := range MemPods {
			if utils.LabelMatch(serviceInfo.Selector, pod.Meta.Labels) {
				podSet = append(podSet, pod)
			}
		}
		//2: 分配信息存储
		for _, pod := range podSet {
			AddPodtoService(pod.Meta.Name, serviceName)
		}
	}
	return err
}

// endregion

// region 删

func deleteService(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		var keys []string
		err = json.Unmarshal(event.Kv.Value, &keys)
		if err != nil {
			return err
		}
		RemoveServices(keys)
		apiserver.ActDeleteService(keys)

	}
	return err
}

//CmdStopPods

// endregion

// region 改

// endregion

// region 查

// endregion

// region 送

// endregion
