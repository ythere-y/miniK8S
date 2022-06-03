package service

import (
	"fmt"
	"io/ioutil"
	"minik8s/meta"
	"minik8s/utils"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

// Container meta data in yaml

type ServiceType string

const (
	ServiceTypeClusterIP    ServiceType = "ClusterIP"
	ServiceTypeNodePort     ServiceType = "NodePort"
	ServiceTypeLoadBalancer ServiceType = "LoadBalancer"
	ServiceTypeExternalName ServiceType = "ExternalName"
)

type ServiceYaml struct {
	Kind     string `yaml:"kind"`
	MetaData struct {
		Name   string            `yaml:"name"`
		Labels map[string]string `yaml:"labels"`
	}
	Selector map[string]string `yaml:"selector"`

	ServiceIP   string `yaml:"serviceip"`
	ServicePort string `yaml:"servicePort"`
	TargetPort  string `yaml:"targetPort"`
}

type Service struct {
	meta.TypeMeta
	meta.ObjectMeat

	Selector    map[string]string
	ServiceIP   string
	ServicePort string
	TargetPort  string
	PodIp       []string
	Type        ServiceType
}

func DisplayService(s Service) {
	fmt.Printf("%-"+strconv.Itoa(blockSize)+"s", s.Name)
	fmt.Printf("%-"+strconv.Itoa(blockSize)+"v", s.Type)
	fmt.Printf("%-"+strconv.Itoa(blockSize)+"v", s.ServiceIP)
	fmt.Printf("%-"+strconv.Itoa(blockSize)+"v", s.ServicePort)

	fmt.Printf("%v", utils.GetAge(s.CreationTimestamp))
	fmt.Println()
}

//
//func (mini *Service) Build(yaml ServiceYaml) {
//	mini.Kind = yaml.Kind
//	mini.Name = yaml.MetaData.Name
//	mini.UID = utils.HashToUid(mini.Name)
//	mini.CreationTimestamp = time.Now()
//
//	mini.ServiceIP = yaml.ServiceIP
//	mini.ServicePort = yaml.ServicePort
//	mini.TargetPort = yaml.TargetPort
//
//	mini.Selector = yaml.Selector
//
//	mini.Type = ServiceTypeClusterIP
//}

var blockSize = 20

func ServicePreDisplay() {

	fmt.Printf("%-"+strconv.Itoa(blockSize)+"s", "NAME")
	fmt.Printf("%-"+strconv.Itoa(blockSize)+"s", "TYPE")
	fmt.Printf("%-"+strconv.Itoa(blockSize)+"s", "ADDR")
	fmt.Printf("%-"+strconv.Itoa(blockSize)+"s", "PORT")
	fmt.Printf("%s", "AGE")
	fmt.Println()
}

func (mini Service) Display() {
	fmt.Printf("%-"+strconv.Itoa(blockSize)+"s", mini.Name)
	fmt.Printf("%-"+strconv.Itoa(blockSize)+"v", mini.Type)
	fmt.Printf("%-"+strconv.Itoa(blockSize)+"v", mini.ServiceIP)
	fmt.Printf("%-"+strconv.Itoa(blockSize)+"v", mini.ServicePort)

	fmt.Printf("%v", utils.GetAge(mini.CreationTimestamp))
	fmt.Println()

	//litter.Dump(mini)
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

func ServiceYamlToService(serviceyaml ServiceYaml) Service {
	var newservice Service
	newservice.Kind = serviceyaml.Kind
	newservice.Name = serviceyaml.MetaData.Name
	newservice.UID = utils.HashToUid(newservice.Name)
	newservice.CreationTimestamp = time.Now()

	newservice.ServiceIP = serviceyaml.ServiceIP
	newservice.ServicePort = serviceyaml.ServicePort
	newservice.TargetPort = serviceyaml.TargetPort

	newservice.Selector = serviceyaml.Selector

	newservice.Type = ServiceTypeClusterIP

	return newservice
}
