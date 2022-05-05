package service

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Container meta data in yaml

type ServiceYaml struct {
	Kind     string `yaml:"kind"`
	MetaData struct {
		Name string `yaml:"name"`
	}
	Selector []string `yaml:"selector,flow"`

	Port       int `yaml:"port"`
	TargetPort int `yaml:"targetPort"`
}

type MiniService struct {
	Kind     string
	MetaData struct {
		Name string
	}
	Selector map[string]string

	Port       int
	TargetPort int
}

func (mini *MiniService) Build(yaml ServiceYaml) {
	mini.Kind = yaml.Kind
	mini.MetaData.Name = yaml.MetaData.Name
	mini.Port = yaml.Port
	mini.TargetPort = yaml.TargetPort
	tmp := make(map[string]string, len(yaml.Selector))
	mini.Selector = tmp
	for _, str := range yaml.Selector {
		//TODO:这部分的转换需要考虑一下
		//fmt.Println(str)
		mini.Selector[str] = "hello_but_nil"
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
