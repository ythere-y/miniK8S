package kubelet

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/client"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"minik8s/apiserver"
	"minik8s/constant"
	"minik8s/environment"
	"minik8s/lab/dksdk"
	. "minik8s/lab/etcd"
	"minik8s/registry/pod"
	"minik8s/utils"
)

// region 听

func kubeletePodWatch() {
	fmt.Printf("kubelet start watch node name = %v\n", NodeName)
	watchName := SetKey(
		SetPrefix(constant.KubeletPrefix),
		SetSourceType(constant.NodeSourceName),
		SetNodeName(NodeName),
	)
	apiserver.SyncWatch(watchName, nodehandler)

}
func nodehandler(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		fmt.Printf("kubelet handling key = %v, value = %v\n", string(event.Kv.Key), string(event.Kv.Value))
		// 增加一个pod的操作
		operation := utils.GetLastWord(string(event.Kv.Key))
		var podInfoGround pod.Pod
		var podInfo *pod.Pod
		err = json.Unmarshal(event.Kv.Value, &podInfoGround)
		podInfo = &podInfoGround
		if err != nil {
			panic(err)
		}
		podName := podInfo.Meta.Name

		switch operation {
		case constant.CREATE:
			// 是创建操作
			CreateAndRunPod(podInfo)
			podIP := environment.GetPodIPByName(podName)
			podInfo.Addr = podIP
		case constant.DELETE:
			// 是删除命令
			StopPod(podName)
			RemovePod(podName)
		default:
			fmt.Printf("op = %v, it not in any!\n", operation)
		}
	}
	return err
}

// endregion

// region 增

func CreateAndRunPod(pod *pod.Pod) uint32 {
	KPods = append(KPods, *pod)
	// currently, use local machine as client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	// create pod with local client

	// TODO: add root container(pause container)
	network := dksdk.RunRootContainer(pod.Meta.Name + "-pause")
	uid := CliCreatePodByPod(cli, *pod, network)
	// uid := CliCreatePodByPod(cli, *pod)
	utils.CheckNil("podsInfo", PodsInfo)
	utils.DebugTanInfo()
	PodsInfo[pod.Meta.Name] = *pod
	for index, ipod := range KPods {
		if ipod.Meta.Uid == uid {
			// allocate cli
			KPods[index].PodClient = cli
		}
	}
	RunPod(uid)
	//RunPodByName(pod.Meta.Name)

	return uid

}

func CliCreatePodByPod(cli *client.Client, pod pod.Pod, net string) uint32 {
	// create meta datas
	newPod := pod
	// allocate client
	newPod.PodClient = cli
	// create containers
	for index, cont := range newPod.Containers {
		image := cont.ContainerImage
		cmd := cont.Command
		var resource dksdk.Resource
		resource.CPUShares = cont.CpuNum
		resource.Memory = cont.Memory
		name := cont.Name
		volumn := cont.Volumn
		port := cont.Port
		hostport := cont.HostPort
		cid := dksdk.CreateContainer(cli, image, cmd,
			resource, name, volumn, port, hostport, net)
		// allocate container id
		newPod.Containers[index].Id = cid
	}
	return newPod.Meta.Uid
}

// endregion

// region 删

//StopPod
/**
 * API: stop pod
 * cli: docker client; podId: pod id specified to stop
**/
func StopPod(name string) {
	for index, pod := range KPods {
		if pod.Meta.Name == name {
			// get its client
			cli := pod.PodClient
			for _, cont := range pod.Containers {
				dksdk.StopContainer(cont.Id, cli)
			}
			// stop manually, failed
			KPods[index].Stats.Status = POD_FAILED
		}
	}
}

//RemovePod
/**
 * API: remove pod
 * cli: docker client; podId: pod to remove
**/
func RemovePod(name string) {
	for index, pod := range KPods {
		if pod.Meta.Name == name {
			// get its client
			cli := pod.PodClient
			// remove containers first
			for _, cont := range pod.Containers {
				dksdk.RemoveContainer(cont.Id, cli)
			}
			// delete pod info from global list
			KPods = append(KPods[:index], KPods[index+1:]...)
		}
	}
}

// endregion

// region 改

// endregion

// region 查

// endregion

// region 令

func RunPod(podId uint32) {
	for index, ipod := range KPods {
		// get specified pod
		if ipod.Meta.Uid == podId {
			// get its client
			cli := ipod.PodClient
			KPods[index].Stats.Status = POD_RUNNING
			// run containers
			for _, cont := range ipod.Containers {
				// start container
				dksdk.StartContainer(cont.Id, cli)
			}
			apiserver.SetPodsStatus(ipod.Meta.Name, KPods[index])
		}
	}
}
func RunPodByName(podName string) {
	ipod := PodsInfo[podName]
	utils.CheckNil("ipod", ipod)
	fmt.Println("pod name = " + ipod.Meta.Name)
	cli := ipod.PodClient
	utils.CheckNil("cli", cli)
	utils.DebugTanInfo()
	fmt.Printf("len = %v\n", len(ipod.Containers))
	for _, cont := range ipod.Containers {
		dksdk.StartContainer(cont.Id, cli)
	}
	utils.DebugTanInfo()
	ipod.Stats.Status = POD_RUNNING
	utils.DebugTanInfo()
	apiserver.SetPodsStatus(ipod.Meta.Name, ipod)
	utils.DebugTanInfo()
	for index, ipod := range KPods {
		utils.DebugTanInfo()
		// get specified pod
		if ipod.Meta.Name == podName {
			KPods[index].Stats.Status = POD_RUNNING
		}
	}
	utils.DebugTanInfo()

}

// endregion
