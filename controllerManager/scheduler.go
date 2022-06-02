package controllerManager

import (
	"fmt"
	"minik8s/apiserver"
	"minik8s/utils"
)

var index = 0

func SchedulerFunc1() string {
	if len(MemNodes) == 0 {
		fmt.Print("no node ~!")
		return ""
	}
	if index >= len(MemNodes) {
		index = 0
	}

	node := MemNodes[index].Name
	index++
	return node
}

func DisPodtoNode(podname string) error {

	var (
		err      error
		nodename string
	)

	nodename = SchedulerFunc1()
	if nodename == "" {
		fmt.Printf("scheduler failed~!")
		return err
	}
	err = apiserver.DistributePodtoNode(nodename, podname)
	utils.HandleError("Distribute pod to node error", err)
	return err

}
