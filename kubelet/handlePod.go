package kubelet

import (
	"encoding/json"
	"fmt"
	"minik8s/apiserver"
	"minik8s/constant"
	"minik8s/lab/dksdk"
	"minik8s/lab/etcd"
	. "minik8s/lab/etcd"
	"minik8s/registry/pod"
	"minik8s/utils"
	"time"

	"github.com/docker/docker/client"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// region 听

func kubeletePodWatch() {

	fmt.Println("[Kubelet] [name = " + NodeName + "] Main started!")

	watchName := SetKey(
		SetPrefix(constant.RelationPrefix),
		SetSourceType(constant.NodeSourceName),
		SetNodeName(NodeName),
	)
	apiserver.SyncWatch(watchName, nodehandler)

	watchName = SetKey(
		SetPrefix(constant.WatchPrefix),
		SetSourceType(constant.NodeSourceName),
		SetNodeName(NodeName),
	)
	apiserver.SyncWatch(watchName, watchHandler)

}
func nodehandler(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.DELETE:
		// 删除一个pod的操作
		fmt.Printf("kubelet handling Delete key = %v\n", string(event.Kv.Key))

		podName := utils.GetLastWord(string(event.Kv.Key))
		cli := StopPod(podName)
		RemovePod(podName)
		dksdk.StopContainer(podName+"-pause", cli)
		dksdk.RemoveContainer(podName+"-pause", cli)
	case mvccpb.PUT:
		fmt.Printf("kubelet handling key = %v, value = %v\n", string(event.Kv.Key), string(event.Kv.Value))
		// 增加/修改 一个pod的操作
		podName := utils.GetLastWord(string(event.Kv.Key))
		operation := string(event.Kv.Value)
		switch operation {
		case podName:
			// 是创建操作
			var podInfo *pod.Pod
			podInfo = apiserver.GetPodInfo(podName)
			CreateAndRunPod(podInfo)
		case constant.StopFlag:
			// 是停止命令
			StopPod(podName)
		case constant.RemoveFlag:
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
	err := apiserver.SavePodInfo(newPod)
	utils.HandleError("save pod info error", err)
	return newPod.Meta.Uid
}

// endregion

// region 删

//StopPod
/**
 * API: stop pod
 * cli: docker client; podId: pod id specified to stop
**/
func StopPod(name string) *client.Client {
	for index, pod := range KPods {
		if pod.Meta.Name == name {
			// get its client
			cli := pod.PodClient
			for _, cont := range pod.Containers {
				dksdk.StopContainer(cont.Id, cli)
			}
			// stop manually, failed
			KPods[index].Stats.Status = POD_FAILED
			return cli
		}
	}
	return nil
}

//RemovePod
/**
 * API: remove pod
 * cli: docker client; podId: pod to remove
**/
func RemovePod(name string) *client.Client {
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
			return cli
		}
	}
	return nil
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

func watchHandler(event *clientv3.Event) error {
	var err error
	switch event.Type {
	case mvccpb.PUT:
		var podName string
		podName = string(event.Kv.Value)
		go kubeletWatchPod(podName)
	}
	return err
}

func kubeletWatchPod(podname string) {
	fmt.Printf("kubelet start to watch pod: %s\n", podname)
	//休眠15秒等pod创建完成
	// time.Sleep(15 * time.Second)
	// fmt.Println("kubelet watch finish sleep")
	cli, err := client.NewClientWithOpts(client.FromEnv)
	utils.HandleError("kubelete watch pod create client error", err)
	//拿到对应的pod
	value, err := etcd.GetValue(
		etcd.SetKey(
			etcd.SetPrefix(constant.RegistryPrefix),
			etcd.SetSourceType(constant.PodSourceName),
			etcd.JustAppend(podname)))
	utils.HandleError("kubelet watch pod get pod error", err)

	var podtmp pod.Pod
	err = json.Unmarshal([]byte(value), &podtmp)
	utils.HandleError("kubelet watch get pod json unmarshal error", err)

	//获取pod里面的容器信息
	var contIds []string
	for _, cont := range podtmp.Containers {
		contid := cont.Id
		contIds = append(contIds, contid)
	}

	go func(cli *client.Client, ids []string, podname string) {
		for {
			//每4秒检查一次
			t := time.NewTicker(4 * time.Second)
			select {
			case <-t.C:
				//检查容器运行情况，如果有容器fail了，就通知master
				for _, id := range ids {
					if !dksdk.IsRun(cli, id) {
						buildkey := etcd.SetKey(
							etcd.SetPrefix(constant.WatchPrefix),
							etcd.SetSourceType(constant.PodSourceName),
							etcd.JustAppend(podname),
						)
						buildvalue := podname
						apiserver.SyncPut(buildkey, buildvalue)
						break
					}
				}
			}
		}
	}(cli, contIds, podname)
}

// endregion
