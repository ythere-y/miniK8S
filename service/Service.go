package service

import (
	"minik8s/apimachinery/pkg/apis/core"
	"minik8s/apimachinery/pkg/apis/meta"
)

type Service struct {
	meta.TypeMeta   `json:",inline" :"meta_._type_meta"`
	meta.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata" :"meta_._object_meta"`
	Spec            core.ServiceSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec" :"spec"`
	Status          core.ServiceStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status" :"status"`
}
