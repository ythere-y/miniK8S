package meta

import "time"

type TypeMeta struct {
	Kind string `yaml:"kind"`
}

type ObjectMeta struct {
	Name               string
	UID                uint32
	CreationTimeStampe time.Time
	DeletionTimeStampe time.Time
	Labels             map[string]string
}
