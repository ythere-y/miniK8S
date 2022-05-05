package service

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"minik8s/meta"
	"minik8s/src/pod"
	"minik8s/utils"
	"time"
)

// Container meta data in yaml

type ServiceYaml struct {
	Kind     string `yaml:"kind"`
	MetaData struct {
		Name   string            `yaml:"name"`
		Labels map[string]string `yaml:"labels"`
	}
	Selector   map[string]string `yaml:"selector"`
	Port       int               `yaml:"port"`
	TargetPort int               `yaml:"targetPort"`
}

type MiniService struct {
	meta.TypeMeta
	meta.ObjectMeta

	Selector   map[string]string
	Pods       []pod.Pod
	Port       int
	TargetPort int
}

func (mini *MiniService) Build(yaml ServiceYaml) {
	mini.TypeMeta.Kind = yaml.Kind
	mini.ObjectMeta.Name = yaml.MetaData.Name
	mini.ObjectMeta.UID = utils.HashToUid(yaml.MetaData.Name)
	mini.ObjectMeta.CreationTimeStampe = time.Now()

	mini.Selector = yaml.Selector
	mini.Port = yaml.Port
	mini.TargetPort = yaml.TargetPort

	// TODO:目前测试将所有Pods加入
	for _, pod := range pod.KPods {
		mini.Pods = append(mini.Pods, pod)
	}
}

func (mini MiniService) Display() {
	fmt.Println(mini)
}

func ParseServiceYaml(file string) ServiceYaml {
	fmt.Println("start parsing yaml file")
	var newPodYaml ServiceYaml
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

func ServiceYamlToService(serviceyaml ServiceYaml) MiniService {
	var newservice MiniService
	newservice.Build(serviceyaml)
	return newservice
}
