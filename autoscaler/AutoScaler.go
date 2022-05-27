package autoscaler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ParseAutoScaler(file string) AutoScaler {
	var asYaml AutoScalerYaml
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("autoscaler yaml file parse err.")
	}

	err = json.Unmarshal(yamlFile, &asYaml)
	if err != nil {
		fmt.Println("unmarshal autoscaler error.")
	}

	var newAutoScaler AutoScaler
	newAutoScaler.Kind = asYaml.Kind
	newAutoScaler.Name = asYaml.Name
	newAutoScaler.Workload = asYaml.Workload
	newAutoScaler.Spec.MinReplicas = asYaml.Spec.MinReplicas
	newAutoScaler.Spec.MaxReplicas = asYaml.Spec.MaxReplicas
	newAutoScaler.Spec.Metrics.Cpu = asYaml.Spec.Metrics.Cpu
	newAutoScaler.Spec.Metrics.Mem = asYaml.Spec.Metrics.Mem

	return newAutoScaler
}
