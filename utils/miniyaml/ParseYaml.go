package miniyaml

import (
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"minik8s/miniInterface"
)

func ParseServiceYaml(file string) miniInterface.ServiceYaml {
	fmt.Println("start parsing yaml file")
	var newPodYaml miniInterface.ServiceYaml
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

func ServiceYamlToService(serviceyaml miniInterface.ServiceYaml) miniInterface.MiniService {
	var newservice miniInterface.MiniService
	newservice.Build(serviceyaml)
	return newservice
}
