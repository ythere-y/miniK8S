package etcd

import (
	"minik8s/utils"
	"strings"
)

var Endpoints []string = nil

func TestConnect() {
	Put("//test part/a", "hello world")
	getRes, err := GetNormal("//test part/a")
	utils.HandleError("get operation error ", err)
	if getRes.Count != 1 {
		utils.DebugError("test connect get error")
	}
	if strings.Compare(string(getRes.Kvs[0].Value), "hello world") != 0 {
		utils.DebugError("test connect get error")
	}
}
