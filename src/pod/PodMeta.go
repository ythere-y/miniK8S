package pod

import (
	"time"

	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
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
	Id             string // allocate when created
	Name           string
	ContainerImage string
	Command        []string
	CpuNum         int64
	Memory         int64
	Volumn         map[string]struct{}
	Port           nat.PortSet
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
	// client in which pod lives, allocate when create
	PodClient *client.Client

	Meta       PodMeta
	Stats      PodStatus // Pod status
	Containers []ContainerMeta
}

// Container meta data in yaml
type ContainerYaml struct {
	Name    string              `yaml:"name"`
	Image   string              `yaml:"image"`
	Command []string            `yaml:"command,flow"`
	Cpu     int64               `yaml:"cpu"`
	Memory  int64               `yaml:"memory"`
	Volumn  map[string]struct{} `yaml:"volumn"`
	Port    nat.PortSet         `yaml:"port"`
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
