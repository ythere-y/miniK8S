package apiserver

import (
	"minik8s/constant"
	"minik8s/lab/etcd"
)

func Reload() {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.ControllerPrefix),
		etcd.JustAppend(constant.RELOAD),
	)
	value := "reload"
	etcd.Put(key, string(value))
}
