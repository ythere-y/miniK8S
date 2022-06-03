package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"minik8s/registry/node"
	"sync"
)

//创建一个结构体
type Config struct {
	EtcdIp         string
	MasterIP       string
	MasterNodeName string
	LocalNodeName  string
	ThisNode       node.Node
}

var Configs Config
var once sync.Once

func Main() {
	get, _ := json.Marshal(Configs)
	fmt.Println(string(get))
}

func init() {
	once.Do(func() {
		//从外部的conf.yaml文件读取数据
		data, _ := ioutil.ReadFile("./config/conf.yaml")
		//使用yaml包，把读取到的data格式化后解析到config实例中
		err := yaml.Unmarshal(data, &Configs)
		if err != nil {
			panic("decode error")
		}
	})
}
