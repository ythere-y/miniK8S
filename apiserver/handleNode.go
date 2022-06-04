package apiserver

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"minik8s/constant"
	"minik8s/lab/etcd"
	"minik8s/registry/node"
	"minik8s/utils"
	"time"
)

var noderole = "_node aipserver_ "

// region 增

func CmdCreateNode(filecontext []byte, cur time.Time) {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.ControllerPrefix),
		etcd.JustAppend(constant.NodeSourceName),
		etcd.JustAppend(constant.CREATE),
		etcd.JustAppend(cur.String()))
	value := filecontext
	log.Println(noderole + "cmd create node")
	SerlPut(key, string(value))
}

func SaveNodeInfo(node node.Node) error {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType(constant.NodeSourceName),
		etcd.SetPodName(node.Name))
	value, err := json.Marshal(node)
	log.Println(noderole + "act save node info")
	utils.HandleError("marshal pod error", err)
	SerlPut(key, string(value))
	return err
}

// endregion

// region 删

// TODO: 实现是移动pod的内容，是有问题的，还需要考量

func ActDeleteNode(names []string) {
	var deleteTargets []string
	for _, key := range names {

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

func CmdDeleteNode(names []string) {
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
		etcd.SetSourceType(constant.NodeSourceName),
		etcd.SetPodName(nodeName))
	value = nodeName
	keys = append(keys, key)
	values = append(values, value)

	key = etcd.SetKey(
		etcd.SetPrefix(constant.KubeletPrefix),
		etcd.SetSourceType(constant.NodeSourceName),
		etcd.SetNodeName(nodeName),
		etcd.SetPodName(podName),
		etcd.JustAppend(constant.CREATE))
	podInfo := GetPodInfo(podName)
	getstr, _ := json.Marshal(podInfo)
	value = string(getstr)
	keys = append(keys, key)
	values = append(values, string(value))

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
