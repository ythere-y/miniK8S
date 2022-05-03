package pod

import (
	"time"
)

// global variable for k8s pods
var KPods []Pod

type PodStatusType string

const (
	POD_PENDING   PodStatusType = "PENDING"
	POD_RUNNING   PodStatusType = "RUNNING"
	POD_SUCCEEDED PodStatusType = "SUCCEEDED"
	POD_FAILED    PodStatusType = "FAILED"
	POD_UNKNOWN   PodStatusType = "UNKNOWN"
)

// Containers meta data in pod
type ContainerMeta struct {
	ContainerImage string
	Command        string
	CpuNum         int
	Memory         int
	Volumn         string
	Port           int
}

// Pod meta data for specification
type PodMeta struct {
	Kind string
	Name string
	Uid  uint32
}

// Pod status
type PodStatus struct {
	CreateTime time.Time
	Status     PodStatusType
}

// Pod structure
type Pod struct {
	Meta       PodMeta
	Stats      PodStatus // Pod status
	Containers []ContainerMeta
}

// Container meta data in yaml
type ContainerYaml struct {
	Image   string `yaml:"image"`
	Command string `yaml:"command"`
	Cpu     int    `yaml:"cpu"`
	Memory  int    `yaml:"memory"`
	Volumn  string `yaml:"volumn"`
	Port    int    `yaml:"port"`
}

type PodYaml struct {
	Kind     string `yaml:"kind"`
	MetaData struct {
		Name string `yaml:"name"`
	}
	Spec struct {
		Containers []ContainerYaml `yaml:"containers,flow"`
	}
}
