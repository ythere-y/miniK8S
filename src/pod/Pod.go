package pod

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// using pod name to get pod uid via hash
func HashToUid(s string) uint32 {
	uid := uint32(crc32.ChecksumIEEE([]byte(s)))
	return uid
}

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

// create a new pod according to yaml file
func CreatePod(file string) {
	// parse yaml file
	podYaml := ParseYaml(file)
	// turn to pod
	pod := YamlToPod(podYaml)
	// add it to global pod list
	KPods = append(KPods, pod)
}

// get all pods info
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

// get one pod info by id
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

// delete one pod by id
func DeletePodById(uid uint32) {
	for index, pod := range KPods {
		if pod.Meta.Uid == uid {
			// delete this pod from slice
			KPods = append(KPods[:index], KPods[index+1:]...)
		}
	}
}

func RunPodById(uid uint32) {
	for index, pod := range KPods {
		if pod.Meta.Uid == uid {
			// set pod status to running
			KPods[index].Stats.Status = POD_RUNNING
			for _, cont := range pod.Containers {
				img := cont.ContainerImage
				// run specific container in background
				RunContainerWithImg(img)
			}
		}
	}
}

func RunPodByName(name string) {
	for index, pod := range KPods {
		if pod.Meta.Name == name {
			// set pod status to running
			KPods[index].Stats.Status = POD_RUNNING
			// run containers
			for _, cont := range pod.Containers {
				img := cont.ContainerImage
				// run specific container in background
				RunContainerWithImg(img)
			}
		}
	}
}
