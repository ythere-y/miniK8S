package apiserver

import (
	"fmt"
	"minik8s/lab/etcd"
	"minik8s/utils"
	"time"
)

func ApiServerTest() {
	etcd.Put(
		etcd.SetKey(etcd.JustAppend("test"),
			etcd.JustAppend("api")), "api-test")
	childNum := GetChildNum("/test/")
	fmt.Printf("children num = %d\n", childNum)

	exist := CheckIfExist("/nothing")
	fmt.Printf("check 1 = noting, %v\n", exist)

	exist = CheckIfExist("/test/a")
	fmt.Printf("check 2 = /test/a, %v\n", exist)
}

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
