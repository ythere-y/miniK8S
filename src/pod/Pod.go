package pod

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"minik8s/lab/dksdk"
	"time"

	"github.com/docker/docker/client"
	yaml "gopkg.in/yaml.v2"
)

/**
 * API: create pod using yaml file
 * specify a docker client, using a yaml file to
 * create a pod, and return its uid.
**/
func CliCreatePod(cli *client.Client, file string) uint32 {
	newPod := ForeHeadCreatePod(file)
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
 * API: run pod
 * cli: docker client, podId: the pod id specified to run
**/
func RunPod(cli *client.Client, podId uint32) {
	for _, pod := range KPods {
		// get specified pod
		if pod.Meta.Uid == podId {
			// run containers
			for _, cont := range pod.Containers {
				// start container
				dksdk.StartContainer(cont.Id, cli)
			}
		}
	}
}

/**
 * API: stop pod
 * cli: docker client; podId: pod id specified to stop
**/
func StopPod(cli *client.Client, podId uint32) {
	for index, pod := range KPods {
		if pod.Meta.Uid == podId {
			for _, cont := range pod.Containers {
				dksdk.StopContainer(cont.Id, cli)
			}
			// stop manually, failed
			KPods[index].Stats.Status = POD_FAILED
		}
	}
}

/**
 * API: remove pod
 * cli: docker client; podId: pod to remove
**/
func RemovePod(cli *client.Client, podId uint32) {
	for index, pod := range KPods {
		if pod.Meta.Uid == podId {
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
 * API: get all pod info
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
 * API: get pod info by uid
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

/* not for api use */
// using pod name to get pod uid via hash
func HashToUid(s string) uint32 {
	uid := uint32(crc32.ChecksumIEEE([]byte(s)))
	return uid
}

/* not for api use */
// turn yaml pod to pod structure
func YamlToPod(podYaml PodYaml) Pod {
	fmt.Println("start creating pod")
	var newPod Pod

	// generate metadata and status
	newPod.Meta.Kind = podYaml.Kind
	newPod.Meta.Name = podYaml.MetaData.Name
	newPod.Meta.Uid = HashToUid(podYaml.MetaData.Name)
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
