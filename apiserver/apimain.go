package apiserver

import (
	"fmt"
	"minik8s/lab/etcd"
	"minik8s/utils"
	"time"
)

func ApiServerMain() {
	fmt.Println("hello world! Api ServerMain started!")
	watchName := etcd.SetKey(etcd.SetPrefix("registry"), etcd.SetSourceType("apiserver"))
	etcd.SyncWatch(watchName, apiserverHandler)

	// sleep
	time.Sleep(3 * time.Second)

	etcd.SyncPut(etcd.SetKey(
		etcd.SetPrefix("registry"),
		etcd.SetSourceType("apiserver"),
		etcd.SetNameSpace("default"),
		etcd.SetName("nickname"),
	), "abc")

	utils.HoldPro()
}
