package service

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/intstr"
	"minik8s/K8SClient"
	"minik8s/apimachinery/pkg/apis/core"
	"minik8s/apimachinery/pkg/apis/meta"
)

type ServiceController struct {
}

func (s ServiceController) CreateService() {
	fmt.Println("create Service")
	CreatServiceTest()
}

func SerMain() {

	var serviceController1 ServiceController

	serviceController1.CreateService()

}

var namespace = "hello"

func CreatServiceTest() {
	// 创建 Service 的选项参数
	//serviceCreateCmd.Flags().StringVar(&serviceCreateName, "name", "", "service name")
	ActuallyCreateRun()
}

func UpdateServiceTest() {
	// 更新 Service 的选项参数
	//serviceUpdateCmd.Flags().StringVar(&serviceUpdateName, "name", "", "service name")
	ActuallyUpdateRun()
}

var serviceCreateName string
var serviceUpdateName string

func ActuallyUpdateRun() {

	// 创建 Kubernetes 客户端对象
	k8sClient, err := K8SClient.CreateK8SClient()
	if err != nil {
		fmt.Println("Err:", err)
		return
	}

	// 获取指定的 Service 对象
	getOptions := meta.GetOptions{}
	service, err := k8sClient.Services(namespace).Get(serviceUpdateName, getOptions)
	if err != nil {
		fmt.Println("Err:", err)
		return
	}

	// 设置原有 Service 暴露的端口，增加 9091
	service.Spec.Ports = append(service.Spec.Ports, core.ServicePort{
		Name:       fmt.Sprintf("tcp-9091-9091-%s", serviceUpdateName),
		Port:       9091,
		TargetPort: intstr.FromInt(9091),
		Protocol:   core.ProtocolTCP,
	})

	// 调用 Update 接口创建
	_, err = k8sClient.Services(namespace).Update(service)
	if err != nil {
		fmt.Println("Err:", err)
		return
	}

	fmt.Println("Update service success!")
}

func ActuallyCreateRun() {

	// 创建 Kubernetes 客户端对象
	/*
		k8sClient, err := K8SClient.CreateK8SClient()
		if err != nil {
			fmt.Println("Err:", err)
			return
		}
	*/
	// 新的 Service 的定义
	var newService core.Service
	var newServiceSpec core.ServiceSpec
	// 设置标签选择器
	newServiceSpec.Selector = map[string]string{
		"app": "echo-go",
	}

	// 设置 Service 端口
	newServiceSpec.Ports = []core.ServicePort{
		core.ServicePort{
			Name:       fmt.Sprintf("tcp-9090-9090-%s", serviceCreateName),
			Port:       9090,
			TargetPort: intstr.FromInt(9090),
			Protocol:   core.ProtocolTCP,
		},
	}

	// 设置 ServiceType 为 NodePort
	newServiceSpec.Type = core.ServiceTypeNodePort

	// 设置 Service 的各个参数
	newService.Spec = newServiceSpec
	newService.Name = serviceCreateName
	newService.Namespace = namespace

	//ServiceController.CreateService(newService)
	// 调用 Create 接口创建
	//_, err = k8sClient.Services(namespace).Create(newService)
	//if err != nil {
	//	fmt.Println("Err:", err)
	//	return
	//}

	fmt.Println("Create service success!")

}
