package apiserver

import (
	"fmt"
	"minik8s/lab/etcd"
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

func Main() {
	fmt.Println("[Api Server] Main started!")
	etcd.TestConnect()

}
