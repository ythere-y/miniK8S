package apiserver

import (
	"minik8s/constant"
	"minik8s/lab/etcd"
	"time"
)

func Reload() {
	key := etcd.SetKey(
		etcd.SetPrefix(constant.ControllerPrefix),
		etcd.JustAppend(constant.RELOAD),
		etcd.JustAppend(time.Now().String()),
	)
	value := "reload"
	SyncPut(key, string(value))
}
