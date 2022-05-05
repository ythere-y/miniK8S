package service

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/intstr"
	"minik8s/K8SClient"
	"minik8s/apimachinery/pkg/apis/core"
	"minik8s/apimachinery/pkg/apis/meta"
	"minik8s/src/pod"
)

var serviceController1 = ServiceController{}

type ServiceController struct {
	ServiceList []MiniService
}

func (s ServiceController) CreateService() {
	fmt.Println("create Service default")

}

func SerReadTest() {
	file := "servicetest.yaml"
	pod.PutPods()

	pod.PodTest()

	BuildService(file)
	GetAllServicInfo()
}

func SerStartCreat() {
	CreateServiceTest()
}

// 打开一个yaml文件构建service并放到servicecontroller的记录中
func BuildService(file string) {
	serviceYaml := ParseServiceYaml(file)
	fmt.Println(serviceYaml)
	service := ServiceYamlToService(serviceYaml)
	serviceController1.ServiceList = append(serviceController1.ServiceList, service)
}

// 打印所有的记录的service的信息
func GetAllServicInfo() {
	for _, service := range serviceController1.ServiceList {
		service.Display()
	}
}

func CreateServiceTest() {
	fmt.Println("create service")

	serviceController1.CreateService()

}

func RemoveAllService() {
	for _, service := range serviceController1.ServiceList {
		for _, pd := range service.Pods {
			pod.RemovePod(pd.Meta.Uid)
		}
	}
}

func DeleteServiceByName(name string) {
	for index, service := range serviceController1.ServiceList {
		if service.Name == name {
			service.DeleteServcie()
		}
		serviceController1.ServiceList = append(serviceController1.ServiceList[:index], serviceController1.ServiceList[index+1:]...)
	}
}
func DeleteServiceByUID(UID uint32) {
	for index, service := range serviceController1.ServiceList {
		if service.UID == UID {
			service.DeleteServcie()
		}
		serviceController1.ServiceList = append(serviceController1.ServiceList[:index], serviceController1.ServiceList[index+1:]...)
	}
}

var namespace = "hello"

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

func init() {
}
