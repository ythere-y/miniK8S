package pod

import (
	"fmt"
	"io/ioutil"
	"minik8s/lab/dksdk"
	"minik8s/utils"
	"time"

	"github.com/docker/docker/client"
	yaml "gopkg.in/yaml.v2"
)

/**
 * API: create pod using yaml file, use default client
 * file: yaml file specify pod structure
**/
func CreatePod(file string) uint32 {
	// currently, use local machine as client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	// create pod with local client
	uid := CliCreatePod(cli, file)
	for index, pod := range KPods {
		if pod.Meta.Uid == uid {
			// allocate cli
			KPods[index].PodClient = cli
		}
	}
	return uid
}

/**
 * API: create pod in a specified client
 * specify a docker client, using a yaml file to
 * create a pod, and return its uid.
**/
func CliCreatePod(cli *client.Client, file string) uint32 {
	// create meta datas
	newPod := ForeHeadCreatePod(file)
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

/**
 * API: run pod on specified client
 * cli: docker client, podId: the pod id specified to run
**/
func RunPod(podId uint32) {
	for index, pod := range KPods {
		// get specified pod
		if pod.Meta.Uid == podId {
			// get its client
			cli := pod.PodClient
			// run containers
			for _, cont := range pod.Containers {
				// start container
				dksdk.StartContainer(cont.Id, cli)
			}
			KPods[index].Stats.Status = POD_RUNNING
		}
	}
}

/**
 * API: run pod on specified client
 * cli: docker client, podId: the pod id specified to run
**/
// func CliRunPod(cli *client.Client, podId uint32) {
// 	for _, pod := range KPods {
// 		// get specified pod
// 		if pod.Meta.Uid == podId {
// 			// run containers
// 			for _, cont := range pod.Containers {
// 				// start container
// 				dksdk.StartContainer(cont.Id, cli)
// 			}
// 		}
// 	}
// }

/**
 * API: stop pod
 * cli: docker client; podId: pod id specified to stop
**/
func StopPod(podId uint32) {
	for index, pod := range KPods {
		if pod.Meta.Uid == podId {
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

/**
 * API: stop pod on a specified client
 * cli: docker client; podId: pod id specified to stop
**/
// func CliStopPod(cli *client.Client, podId uint32) {
// 	for index, pod := range KPods {
// 		if pod.Meta.Uid == podId {
// 			for _, cont := range pod.Containers {
// 				dksdk.StopContainer(cont.Id, cli)
// 			}
// 			// stop manually, failed
// 			KPods[index].Stats.Status = POD_FAILED
// 		}
// 	}
// }

/**
 * API: remove pod
 * cli: docker client; podId: pod to remove
**/
func RemovePod(podId uint32) {
	for index, pod := range KPods {
		if pod.Meta.Uid == podId {
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

/**
 * API: remove pod on a specified client
 * cli: docker client; podId: pod to remove
**/
// func CliRemovePod(cli *client.Client, podId uint32) {
// 	for index, pod := range KPods {
// 		if pod.Meta.Uid == podId {
// 			// remove containers first
// 			for _, cont := range pod.Containers {
// 				dksdk.RemoveContainer(cont.Id, cli)
// 			}
// 			// delete pod info from global list
// 			KPods = append(KPods[:index], KPods[index+1:]...)
// 		}
// 	}
// }

/**
 * API: print all pod info
**/
func GetAllPodInfo() {
	fmt.Println("UID, name, status, living_time")
	for _, pod := range KPods {
		uid := pod.Meta.Uid
		name := pod.Meta.Name
		status := pod.Stats.Status
		nowtime := time.Now()
		duration := nowtime.Sub(pod.Stats.CreateTime)
		livingTime := duration.String()
		fmt.Printf("%d, %s, %s, %s\n", uid, name, status, livingTime)
	}
}

/**
 * API: print pod info by uid
 * uid: pod id
**/
func GetPodInfoById(uid uint32) {
	for _, pod := range KPods {
		if pod.Meta.Uid == uid {
			uid := pod.Meta.Uid
			name := pod.Meta.Name
			status := pod.Stats.Status
			nowtime := time.Now()
			duration := nowtime.Sub(pod.Stats.CreateTime)
			livingTime := duration.String()
			fmt.Printf("%d, %s, %s, %s\n", uid, name, status, livingTime)
		}
	}
}

/**
 * API: print containers info in specified pod
 * uid: pod id
**/
func PrintContainerInfoInPod(uid uint32) {
	for _, pod := range KPods {
		// find the pod
		if pod.Meta.Uid == uid {
			// get the client
			cli := pod.PodClient
			for _, cont := range pod.Containers {
				// print container info
				dksdk.ContainerStat(cli, cont.Id)
			}
		}
	}
}

/* not for api use */
// turn yaml pod to pod structure
func YamlToPod(podYaml PodYaml) Pod {
	fmt.Println("start creating pod")
	var newPod Pod

	// generate metadata and status
	newPod.Meta.Kind = podYaml.Kind
	newPod.Meta.Name = podYaml.MetaData.Name
	newPod.Meta.Uid = utils.HashToUid(podYaml.MetaData.Name)
	newPod.Meta.Labels = podYaml.MetaData.Labels
	newPod.Stats.CreateTime = time.Now()
	newPod.Stats.Status = POD_PENDING

	// deal with containers
	var tmpContainer ContainerMeta
	for _, value := range podYaml.Spec.Containers {
		tmpContainer.Id = ""
		tmpContainer.ContainerImage = value.Image
		tmpContainer.Command = value.Command
		tmpContainer.CpuNum = value.Cpu
		tmpContainer.Memory = value.Memory
		tmpContainer.Volumn = value.Volumn
		tmpContainer.Port = value.Port
		// append it to pod
		newPod.Containers = append(newPod.Containers, tmpContainer)
	}

	return newPod
}

/* not for api use */
// parse yaml file to a yaml pod structure
func ParseYaml(file string) PodYaml {
	fmt.Println("start parsing yaml file")
	var newPodYaml PodYaml
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("yaml file read error")
	}

	err = yaml.Unmarshal(yamlFile, &newPodYaml)
	if err != nil {
		fmt.Println("yaml unmarshal error")
	}
	return newPodYaml
}

/* not for api use */
// create a new pod to global var according to yaml file,
func ForeHeadCreatePod(file string) Pod {
	// parse yaml file
	podYaml := ParseYaml(file)
	// turn to pod
	pod := YamlToPod(podYaml)
	// add it to global pod list
	KPods = append(KPods, pod)
	for index, value := range KPods {
		if pod.Meta.Uid == value.Meta.Uid {
			return KPods[index]
		}
	}
	// should not be here
	return KPods[0]
}
