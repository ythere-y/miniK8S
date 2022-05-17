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
	go etcd.SyncWatch("t1")
	go etcd.SyncWatch("/registry/test/default/")
	//go etcd.SyncWatch("/registry/test/default/nickname")
	//go etcd.SyncWatch("/registry/test/default")
	time.Sleep(3 * time.Second)
	go etcd.Put("registry", "test", "default", "nickname", "abc")
	//go etcd.SyncPutTest("t1", "hello world")
	//go etcd.SyncPutTest("t1", "hello")
	//go etcd.SyncPutTest("t1", "world")
	utils.HoldPro()
}
