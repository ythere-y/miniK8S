package apiserver

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"minik8s/constant"
	"minik8s/lab/etcd"
	"minik8s/registry/node"
	"minik8s/utils"
	"time"
)

// region 增

func CreateNode(status node.Node) {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType(constant.NodeSourceName),
		etcd.SetNodeName(status.Name))
	value, err := json.Marshal(status)
	utils.HandleError("marshal node status error", err)
	SyncPut(key, string(value))
}

// endregion

// region 删

func DeleteNode(names []string) {
	// 检查是否存在在relation关系中
	nameSet := utils.ParseNames(names)
	for _, name := range nameSet {
		path := etcd.SetKey(
			etcd.SetPrefix(constant.RelationPrefix),
			etcd.SetSourceType(constant.NodeSourceName),
			etcd.SetPodName(name))
		if CheckIfExist(path) == false {
			fmt.Println("cannot find node [" + name + "] !")
			return
		}
	}
	names = nameSet

	// 准备给controller发送命令
	key := etcd.SetKey(
		etcd.SetPrefix(constant.ControllerPrefix),
		etcd.SetSourceType(constant.NodeSourceName),
		etcd.JustAppend(constant.DELETE),
		etcd.JustAppend(time.Now().String()))
	value, err := json.Marshal(names)
	if err != nil {
		panic(err)
		return
	}
	SyncPut(key, string(value))
}

// endregion

// region 改

func SetNodeStatus(name string, target node.Node) {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType(constant.PodSourceName),
		etcd.SetPodName(name))
	value, _ := json.Marshal(target)
	etcd.Put(key, string(value))
}

//DistributePodtoNode
/*
将一个pod分配给一个node（两者使用name来识别）
此为执行操作，将分配结果写入etcd
*/
func DistributePodtoNode(nodeName string, podName string) error {
	var (
		key    string
		value  string
		keys   []string
		values []string
	)

	key = etcd.SetKey(
		etcd.SetPrefix(constant.RelationPrefix),
		etcd.SetSourceType(constant.NodeSourceName),
		etcd.SetNodeName(nodeName),
		etcd.SetPodName(podName))
	value = podName
	keys = append(keys, key)
	values = append(values, value)

	key = etcd.SetKey(
		etcd.SetPrefix(constant.RelationPrefix),
		etcd.SetSourceType(constant.PodSourceName),
		etcd.SetNodeName(podName),
		etcd.SetPodName(nodeName))
	value = nodeName
	keys = append(keys, key)
	values = append(values, value)

	SyncPutList(keys, values)

	return nil
}

// endregion

// region 查

//GetNodeInfo
/*
 根据podName，从/registry/pods/目录下寻找对应的pod
并组建成Pod返回
*/
func GetNodeInfo(name string) *node.Node {
	getRes, err := etcd.Get(etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType(constant.NodeSourceName),
		etcd.SetName(name)))

	utils.HandleError("get node info error", err)
	if len(getRes) == 0 {
		return nil
	}
	var nodeInfo node.Node
	err = json.Unmarshal([]byte(getRes[0]), &nodeInfo)
	utils.HandleError("unmarshal node failed ", err)

	return &nodeInfo
}

func DisplayAllNodeInfo() {
	var (
		num    uint32 = 0
		getRsp *clientv3.GetResponse
		err    error
	)
	getRsp, err = etcd.GetWithPrefix(
		etcd.SetKey(
			etcd.SetPrefix(constant.RegistryPrefix),
			etcd.SetSourceType(constant.NodeSourceName)))
	utils.HandleError("get with prefix error[from get child num]", err)
	num = uint32(len(getRsp.Kvs))
	if num == 0 {
		fmt.Println("cannot find any nodes")
		return
	}
	node.NodePreDisplay()
	for _, event := range getRsp.Kvs {
		tmpValue := event.Value
		var tmpPod node.Node
		err = json.Unmarshal(tmpValue, &tmpPod)
		utils.HandleError("unmarshal node error", err)
		tmpPod.Display()
	}
}

//DisplayNodeInfo
/*
按照podName查找Pods的信息（以表格形式打印主要信息）
*/
func DisplayNodeInfo(name string) {
	tmpNode := GetNodeInfo(name)
	if tmpNode == nil {
		fmt.Println("cannot find any node")
	} else {
		node.NodePreDisplay()
		tmpNode.Display()
	}
}

// endregion
