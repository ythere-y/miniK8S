package controllerManager

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"minik8s/apiserver"
	"minik8s/config"
	"minik8s/constant"
	. "minik8s/lab/etcd"
	"minik8s/registry/node"
	"minik8s/utils"
)

var noderole = " __node controller__ "

func nodeControllerWatch() {

	fmt.Println("[Node controller] init!")
	var watchName string
	watchName = SetKey(
		SetPrefix(constant.ControllerPrefix),
		SetSourceType(constant.NodeSourceName),
		JustAppend(constant.CREATE),
	)
	apiserver.SyncWatch(watchName, createNode)

	watchName = SetKey(
		SetPrefix(constant.ControllerPrefix),
		SetSourceType(constant.NodeSourceName),
		JustAppend(constant.DELETE),
	)
	apiserver.SyncWatch(watchName, deleteNode)

}

// region 增

func CreateMasterNode(nodefile []byte) {
	nodeInfo := node.NodeYamlToNode(node.ParseNodeYaml(nodefile))
	AddNode(nodeInfo)
	config.SetConfigThisNode(nodeInfo)
	displayCurMemNodes()
}
func createNode(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		log.Println(noderole + "create Node")
		var newone node.NodeYaml
		newone = node.ParseNodeYaml(event.Kv.Value)
		if err != nil {
			fmt.Printf("yaml unmarshal error->:\n%v\n", err.Error())
		}
		nodeInfo := node.NodeYamlToNode(newone)
		nodename := nodeInfo.Name
		log.Printf("node name = %v\n", nodename)
		for _, memNode := range MemNodes {
			if memNode.Name == nodename {
				fmt.Printf("node %v already exist~!\n", nodename)
				apiserver.Reply(event.Kv.Key, constant.ReplayERROR)
				return err
			}
		}

		AddNode(nodeInfo)

		apiserver.Reply(event.Kv.Key, constant.ReplayOK)
		err = apiserver.SaveNodeInfo(nodeInfo)

		utils.HandleError("save pod info error", err)
		displayCurMemNodes()
	}
	return err
}

// endregion

// region 删

func deleteNode(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		var keys []string
		err = json.Unmarshal(event.Kv.Value, &keys)
		if err != nil {
			return err
		}

		var podSet []string
		for _, nodename := range keys {
			for _, podname := range relations.NodetoPodRela[nodename] {
				podSet = append(podSet, podname)
			}
		}
		//先关闭所有的pod，再关闭所有的node
		RemovePods(podSet)
		apiserver.ActDeletePods(podSet)
		RemoveNodes(keys)
		apiserver.ActDeleteNode(keys)
	}
	return err

}

// endregion

// region 改

// endregion

// region 查

func displayCurMemNodes() {
	fmt.Printf("after node create , memnodes display\n")
	for i, memNode := range MemNodes {
		js, _ := json.Marshal(memNode)
		fmt.Printf("[node %v] = %v\n", i, string(js))
	}
}

// endregion

// region 送

// endregion
