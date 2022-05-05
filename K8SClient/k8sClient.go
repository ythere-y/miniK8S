package K8SClient

import (
	"minik8s/apimachinery/pkg/apis/core"
)

//TODO:Client的数据结构

type K8sClient struct {
}

//TODO:在client中获取services

func (client *K8sClient) Services(namespace string) core.Service {
	var service core.Service

	return service
}

//TODO:如何创建一个client

func CreateK8SClient() (*K8sClient, error) {
	var client K8sClient
	var err error

	return &client, err
}
