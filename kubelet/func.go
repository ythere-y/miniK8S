package kubelet

import (
	"fmt"
	"github.com/docker/docker/client"
	"minik8s/apiserver"
	"minik8s/lab/dksdk"
	"minik8s/pod"
	"minik8s/utils"
)

var (
	KPods    []pod.Pod
	PodsInfo map[string]pod.Pod
)

func CreateAndRunPod(pod *pod.Pod) uint32 {
	KPods = append(KPods, *pod)
	// currently, use local machine as client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	// create pod with local client

	uid := CliCreatePodByPod(cli, *pod)
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

func CliCreatePodByPod(cli *client.Client, pod pod.Pod) uint32 {
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
		volume := cont.Volumn
		port := cont.Port
		cid := dksdk.CreateContainer(cli, image, cmd,
			resource, name, volume, port, "")
		// allocate container id
		newPod.Containers[index].Id = cid
	}
	return newPod.Meta.Uid
}
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

func init() {
	PodsInfo = make(map[string]pod.Pod)
}
