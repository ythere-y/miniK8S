package replicaset

import (
	"fmt"
	"io/ioutil"
	"minik8s/pod"
	"minik8s/utils"

	yaml "gopkg.in/yaml.v2"
)

/* not for api use */
// parse rs yaml file to rs yaml struct
func ParseRSYaml(file string) RSyaml {
	fmt.Println("start parsing yaml file to rs!")
	var newRSyaml RSyaml
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("rs yaml file read error!")
	}
	err = yaml.Unmarshal(yamlFile, &newRSyaml)
	if err != nil {
		fmt.Println("rs yaml file unmarshal error!")
	}
	return newRSyaml
}

/* not for api use */
// turn rs yaml struct to replicaset struct
func YamlToRS(rsyaml RSyaml) ReplicaSet {
	var newRS ReplicaSet
	newRS.RSmeta.Kind = rsyaml.Kind
	newRS.RSmeta.Name = rsyaml.Metadata.Name
	newRS.RSmeta.Uid = utils.HashToUid(newRS.RSmeta.Name)
	newRS.RSspec.Replicas = rsyaml.Spec.Replicas
	newRS.RSspec.SelectorLabels = rsyaml.Spec.Selector.MatchLabels
	newRS.PodTemplate.Meta.Name = rsyaml.Spec.Template.Metadata.Name
	newRS.PodTemplate.Meta.Labels = rsyaml.Spec.Template.Metadata.Labels

	var tmpContainer pod.ContainerMeta
	for _, value := range rsyaml.Spec.Template.Spec.Containers {
		tmpContainer.Id = ""
		tmpContainer.ContainerImage = value.Image
		tmpContainer.Command = value.Command
		tmpContainer.CpuNum = value.Cpu
		tmpContainer.Memory = value.Memory
		tmpContainer.Volumn = value.Volumn
		tmpContainer.Port = value.Port
		// append it to pod
		newRS.PodTemplate.Containers = append(newRS.PodTemplate.Containers, tmpContainer)
	}

	return newRS
}

// parse podtemplate in replicaset structure to create a actual pod instance
func CreatePodByRsPodTemplate(rs ReplicaSet) pod.Pod {
	podTemplate := rs.PodTemplate

}

/*
 * API use: parse a yaml file to replicaset structure
 * input: rs yaml file
 * return: a replicaSet structure instance
 */
func ParseYamlToRS(file string) ReplicaSet {
	rsyaml := ParseRSYaml(file)
	newRS := YamlToRS(rsyaml)
	return newRS
}
