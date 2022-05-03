package miniyaml

import (
	"fmt"
	"io/ioutil"
	"minik8s/service"

	yaml "gopkg.in/yaml.v2"
)

func ParseServiceYaml(file string) service.ServiceYaml {
	fmt.Println("start parsing yaml file")
	var newPodYaml service.ServiceYaml
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

func ServiceYamlToService(serviceyaml service.ServiceYaml) service.MiniService {
	var newservice service.MiniService
	newservice.Build(serviceyaml)
	return newservice
}
