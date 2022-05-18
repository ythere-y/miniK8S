package apiserver

import (
	"fmt"
	"minik8s/lab/etcd"
	"minik8s/utils"
	"time"
)

var test_lock_1 = true

var test_lock_2 = true

func apiserverstart() {
	fmt.Println("api server start ok!")
	for test_lock_1 {

	}
	fmt.Println("lock_1 get")
	time.Sleep(10 * time.Second)
	test_lock_2 = false
}

func sendOperation() {
	test_lock_1 = false

	for test_lock_2 {

	}
	fmt.Println("lock_2 get")
	fmt.Println("finished all")
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
