package apiserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"minik8s/constant"
	"minik8s/lab/etcd"
	"minik8s/registry/pod"
	"minik8s/registry/replicaset"
	"minik8s/utils"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

//DisplayAllRsInfo
/*
 * display all pods info in all replicasets
 */
func DisplayAllRsInfo() {
	var (
		num    uint32 = 0
		getRsp *clientv3.GetResponse
		err    error
	)
	getRsp, err = etcd.GetWithPrefix(
		etcd.SetKey(
			etcd.SetPrefix(constant.RegistryPrefix),
			etcd.SetSourceType(constant.ReplicaSourceName)))
	utils.HandleError("get with prefix error[from get child num]", err)
	num = uint32(len(getRsp.Kvs))
	if num == 0 {
		fmt.Println("cannot find any replicaset")
		return
	}
	pod.PodPreDisplay()
	for _, event := range getRsp.Kvs {
		tmpValue := event.Value
		var tmpRs replicaset.ReplicaSet
		err = json.Unmarshal(tmpValue, &tmpRs)
		utils.HandleError("unmarshal rs in display error", err)
		DisplayPodsInRs(tmpRs)
	}
}

//DisplayRsInfo
func DisplayRsInfo(name string) {
	tmpRs := GetRsInfo(name)
	if tmpRs == nil {
		fmt.Println("cannot find this replicaset")
	} else {
		pod.PodPreDisplay()
		DisplayPodsInRs(*tmpRs)
	}
}

// not for api use
func DisplayPodsInRs(rs replicaset.ReplicaSet) {
	replicas := rs.RSspec.Replicas
	podname := rs.PodTemplate.Meta.Name
	for i := 1; i <= replicas; i++ {
		tmpname := podname + "-" + strconv.Itoa(i)
		tmpPod := GetPodInfo(tmpname)
		tmpPod.Display()
	}
}

//GetRsInfo
/*
根据name，获取存储的replicaset
*/
func GetRsInfo(rsName string) *replicaset.ReplicaSet {
	var rsInfo replicaset.ReplicaSet
	getRes, err := etcd.Get(etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType(constant.ReplicaSourceName),
		etcd.JustAppend(rsName)))
	utils.HandleError("get rs info error", err)
	if len(getRes) == 0 {
		return nil
	}
	err = json.Unmarshal([]byte(getRes[0]), &rsInfo)
	utils.HandleError("unmarshal rs info failed", err)
	return &rsInfo
}

// Replicaset related functions
/*
 * create a replicaset
 * input: yaml file
 */
func CreateRs(file string) {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.ControllerPrefix),
		etcd.JustAppend(constant.ReplicaSourceName),
		etcd.JustAppend(constant.CREATE),
		etcd.JustAppend(time.Now().String()))
	value, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("file %v read error!\n", file)
	}
	SyncPut(key, string(value))
}

/*
 * delete replicasets
 * input: replicasets names to be deleted
 */
func DeleteRs(names []string) error {
	nameSet := utils.ParseNames(names)
	for _, name := range nameSet {
		path := etcd.SetKey(
			etcd.SetPrefix(constant.RegistryPrefix),
			etcd.SetSourceType(constant.ReplicaSourceName),
			etcd.SetPodName(name))
		if CheckIfExist(path) == false {
			fmt.Println("cannot find replicaset [" + name + "] !")
			return nil
		}
	}
	names = nameSet

	key := etcd.SetKey(
		etcd.SetPrefix(constant.ControllerPrefix),
		etcd.SetSourceType(constant.ReplicaSourceName),
		etcd.JustAppend(constant.DELETE),
		etcd.JustAppend(time.Now().String()))
	value, err := json.Marshal(names)
	if err != nil {
		return err
	}
	SyncPut(key, string(value))
	return nil
}

// save replicaset info
func SaveRsInfo(rs replicaset.ReplicaSet) error {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType(constant.ReplicaSourceName),
		etcd.SetPodName(rs.RSmeta.Name))
	value, err := json.Marshal(rs)

	utils.HandleError("marshal rs error", err)
	etcd.Put(key, string(value))
	return err
}

/*
 * pods managed by replicaset should be put into a relation
 * part of etcd: /relation/replicaset/[rsname]/[podname]
 */
func SaveRsPodInfo(rs replicaset.ReplicaSet, pod pod.Pod) error {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RelationPrefix),
		etcd.SetSourceType(constant.ReplicaSourceName),
		etcd.JustAppend(rs.RSmeta.Name),
		etcd.JustAppend(pod.Meta.Name))
	value := pod.Meta.Name
	etcd.Put(key, string(value))
	return nil
}

func CheckRs(rs replicaset.ReplicaSet, pod pod.Pod) error {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RelationPrefix),
		etcd.SetSourceType(constant.ReplicaSourceName),
		etcd.JustAppend(rs.RSmeta.Name),
		etcd.JustAppend(pod.Meta.Name))
	value := pod.Meta.Name
	etcd.Put(key, string(value))
	return nil
}
