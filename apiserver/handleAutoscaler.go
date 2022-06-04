package apiserver

import (
	"encoding/json"
	"minik8s/constant"
	"minik8s/lab/etcd"
	"minik8s/registry/autoscaler"
	"minik8s/utils"
)

func CmdCreateAs(file string) {
	newAs := autoscaler.ParseAutoScaler(file)
	SaveAsInfo(newAs)
}

func SaveAsInfo(as autoscaler.AutoScaler) error {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.RegistryPrefix),
		etcd.SetSourceType(constant.AutoscalerSourceName),
		etcd.JustAppend(as.Name),
	)
	value, err := json.Marshal(as)
	utils.HandleError("save as info error", err)
	SyncPut(key, string(value))
	return err
}
