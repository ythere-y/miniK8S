package controllerManager

import (
	"encoding/json"
	"fmt"
	"minik8s/apiserver"
	"minik8s/constant"
	"minik8s/lab/etcd"
	. "minik8s/lab/etcd"
	"minik8s/registry/replicaset"
	"minik8s/utils"
	"reflect"
	"strconv"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v2"
)

func init() {
	fmt.Printf("Pod controller init!")
	var watchName string
	watchName = SetKey(
		SetPrefix(constant.ControllerPrefix),
		SetSourceType(constant.ReplicaSourceName),
		JustAppend(constant.CREATE),
	)
	apiserver.SyncWatch(watchName, CreateRs)

	watchName = SetKey(
		SetPrefix(constant.ControllerPrefix),
		SetSourceType(constant.ReplicaSourceName),
		JustAppend(constant.DELETE),
	)
	apiserver.SyncWatch(watchName, DeleteRs)
}

// region 增

func CreateRs(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		fmt.Printf("controller rs start creating")
		var newRsYaml replicaset.RSyaml
		err = yaml.Unmarshal(event.Kv.Value, &newRsYaml)
		if err != nil {
			fmt.Printf("yaml unmarshal error->:\n%v\n", err.Error())
		}
		rsInfo := replicaset.YamlToRS(newRsYaml)
		err = apiserver.SaveRsInfo(rsInfo)
		utils.HandleError("save replicaset info error", err)
		// TODO: deal with pod replicas, create pod in some nodes. UNFINISHED
		if reflect.DeepEqual(rsInfo.RSspec.SelectorLabels,
			rsInfo.PodTemplate.Meta.Labels) {
			fmt.Printf("controller rs create pods")
			replicas := rsInfo.RSspec.Replicas
			pods := replicaset.CreatePodInstances(rsInfo, replicas)
			for _, pod := range pods {
				err = apiserver.SaveRsPodInfo(rsInfo, pod)
				utils.HandleError("save rs info error", err)
				// create pod
				AddPod(pod)
				err = apiserver.SavePodInfo(pod)
				utils.HandleError("save pod info error", err)

				// 这里取消了scheduler的数据转发层，直接把scheduler作为controller内部的一个分支
				err = DisPodtoNode(pod.Meta.Name)
			}
		}
	}
	return err
}

// endregion

// region 删

// TODO  delete的步骤还需要先stop，还需要考虑如何stop结束之后再delete
func DeleteRs(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		var names []string
		var rsDelTar []string
		var podDelTar []string
		err = json.Unmarshal(event.Kv.Value, &names)
		for _, rsname := range names {
			// set rs key to be deleted
			rskey := etcd.SetKey(
				etcd.SetPrefix(constant.RegistryPrefix),
				etcd.SetSourceType(constant.ReplicaSourceName),
				etcd.JustAppend(rsname))
			rsDelTar = append(rsDelTar, rskey)

			// set pod key to be deleted
			rsInfo := apiserver.GetRsInfo(rsname)
			replicas := rsInfo.RSspec.Replicas
			podname := rsInfo.PodTemplate.Meta.Name
			for i := 1; i <= replicas; i++ {
				rspodname := podname + "-" + strconv.Itoa(i)
				// 之后会交给deletePod处理
				podDelTar = append(podDelTar, rspodname)
			}
		}
		// delete rs info
		apiserver.SyncDel(rsDelTar)
		// delete pods in these rs
		apiserver.CmdStopPods(podDelTar)
		apiserver.CmdDeleteNode(podDelTar)
	}
	return err
}
