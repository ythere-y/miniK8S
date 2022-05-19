package replicaset

import (
	"minik8s/pod"
)

type ReplicaSetMeta struct {
	Kind string
	Name string
	Uid  uint32
}

type ReplicaSetSpec struct {
	Replicas       int
	SelectorLabels map[string]string
}

// ReplicaSet abstract struct
type ReplicaSet struct {
	RSmeta      ReplicaSetMeta
	RSspec      ReplicaSetSpec
	PodTemplate pod.Pod
}

type RSpodyaml struct {
	Metadata struct {
		Name   string            `yaml:"name"`
		Labels map[string]string `yaml:"labels"`
	}
	Spec struct {
		Containers []pod.ContainerYaml `yaml:"containers"`
	}
}

// rs struct for parsing yaml file
type RSyaml struct {
	Kind     string `yaml:"kind"`
	Metadata struct {
		Name string `yaml:"name"`
	}
	Spec struct {
		Replicas int `yaml:"replicas"`
		Selector struct {
			MatchLabels map[string]string `yaml:"matchlabels"`
		}
		Template RSpodyaml `yaml:"template"`
	}
}
