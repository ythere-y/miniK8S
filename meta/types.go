package meta

import "time"

type TypeMeta struct {
	Kind string
}

type ObjectMeat struct {
	Name              string
	UID               uint32
	CreationTimestamp time.Time
	DeletionTimestamp time.Time
}
