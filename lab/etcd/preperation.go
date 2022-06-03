package etcd

import (
	"fmt"
	"minik8s/config"
	"minik8s/utils"
	"strings"
)

var Endpoints []string = nil

func TestConnect() {
	fmt.Printf("try to connect to etcd, etcdIP = %v\n", config.Configs.EtcdIp)
	Put("//test part/a", "hello world")
	getRes, err := GetNormal("//test part/a")
	utils.HandleError("get operation error ", err)
	if getRes.Count != 1 {
		utils.DebugError("test connect get error")
	}
	if strings.Compare(string(getRes.Kvs[0].Value), "hello world") != 0 {
		utils.DebugError("test connect get error")
	}
	fmt.Printf("test connect success!")
}
