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
	Selector map[string]string `yaml:"selector"`

	Port       int `yaml:"port"`
	TargetPort int `yaml:"targetPort"`
}

type MiniService struct {
	meta.TypeMeta
	meta.ObjectMeat

	Selector   map[string]string
	Pods       []pod.Pod
	Port       int
	TargetPort int
}

func (mini *MiniService) Build(yaml ServiceYaml) {
	mini.Kind = yaml.Kind
	mini.Name = yaml.MetaData.Name
	mini.UID = utils.HashToUid(mini.Name)
	mini.CreationTimestamp = time.Now()

	mini.Port = yaml.Port
	mini.TargetPort = yaml.TargetPort
	mini.Selector = yaml.Selector

	for _, pod := range pod.KPods {
		mini.Pods = append(mini.Pods, pod)
	}
}

func (mini MiniService) Display() {
	fmt.Println(mini)
}

type IntOrString struct {
	Type   Type   `protobuf:"varint,1,opt,name=type,casttype=Type"`
	IntVal int32  `protobuf:"varint,2,opt,name=intVal"`
	StrVal string `protobuf:"bytes,3,opt,name=strVal"`
}

// Type represents the stored type of IntOrString.
type Type int64

const (
	Int    Type = iota // The IntOrString holds an int.
	String             // The IntOrString holds a string.
)

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
