package kubelet

import (
	"github.com/docker/docker/client"
	"minik8s/lab/dksdk"
	"minik8s/src/pod"
)

var (
	KPods []pod.Pod
)

func CreateAndRunPod(pod pod.Pod) uint32 {
	KPods = append(KPods, pod)
	// currently, use local machine as client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	// create pod with local client

	uid := CliCreatePodByPod(cli, pod)
	for index, pod := range KPods {
		if pod.Meta.Uid == uid {
			// allocate cli
			KPods[index].PodClient = cli
		}
	}
	RunPod(uid)

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
		var resouce dksdk.Resource
		resouce.CPUShares = cont.CpuNum
		resouce.Memory = cont.Memory
		name := cont.Name
		volume := cont.Volumn
		port := cont.Port
		cid := dksdk.CreateContainer(cli, image, cmd,
			resouce, name, volume, port, "")
		// allocate container id
		newPod.Containers[index].Id = cid
	}
	return newPod.Meta.Uid
}
func RunPod(podId uint32) {
	for index, pod := range KPods {
		// get specified pod
		if pod.Meta.Uid == podId {
			// get its client
			cli := pod.PodClient
			KPods[index].Stats.Status = POD_RUNNING
			// run containers
			for _, cont := range pod.Containers {
				// start container
				dksdk.StartContainer(cont.Id, cli)
			}
		}
	}
}
