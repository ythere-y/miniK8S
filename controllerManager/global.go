package controllerManager

import (
	"minik8s/apiserver"
	"minik8s/registry/node"
	pod2 "minik8s/registry/pod"
	"minik8s/registry/service"
)

func AddPod(status pod2.Pod) {
	MemPods = append(MemPods, status)
}

func DeletePod(name string) {
	len := len(MemPods)
	for i := 0; i < len; i++ {
		if MemPods[i].Meta.Name == name {
			MemPods = append(MemPods[:i], MemPods[i+1:]...)
			break
		}
	}
}
func RemovePods(names string) pod2.Pod {
	for i, pod := range MemPods {
		if pod.Meta.Name == name {
			MemPods = append(MemPods[:i], MemPods[i+1:]...)
			return pod
			break
		}
	}
	return nil
}

func AddNode(node node.Node) {
	MemNodes = append(MemNodes, node)
}
func DeleteNode(name string) {
	len := len(MemNodes)
	for i := 0; i < len; i++ {
		if MemNodes[i].Name == name {
			MemNodes = append(MemNodes[:i], MemNodes[i+1:]...)
			break
		}
	}
}
func RemoveNodes(names []string) {
	for _, name := range names {
		for i, node := range MemNodes {
			if node.Name == name {
				MemNodes = append(MemNodes[:i], MemNodes[i+1:]...)
				break
			}
		}
	}
}

func MemDeleteService(name string) {
	len := len(MemServices)
	for i := 0; i < len; i++ {
		if MemServices[i].Name == name {
			MemServices = append(MemServices[:i], MemServices[i+1:]...)
			break
		}
	}
	delete(relations.ServicetoPodRela, name)
}

func AddService(service service.Service) {
	MemServices = append(MemServices, service)
}
func RemoveServices(name string) service.Service {
	for i, memService := range MemServices {
		if memService.Name == name {
			MemServices = append(MemServices[:i], MemServices[i+1:]...)
			return memService
			break
		}
	}
	return nil
}

func AddPodtoNode(podname string, nodename string) {
	// 从原来的地方移除
	for pod, nodeidx := range relations.PodstoNodeRela {
		if pod == podname {
			for i, po := range relations.NodetoPodRela[nodeidx] {
				if po == podname {
					relations.NodetoPodRela[nodeidx] = append(
						relations.NodetoPodRela[nodeidx][:i],
						relations.NodetoPodRela[nodeidx][i+1:]...,
					)
				}
			}
		}
	}
	relations.PodstoNodeRela[podname] = nodename
	relations.NodetoPodRela[nodename] = append(relations.NodetoPodRela[nodename], podname)
}
func AddPodtoService(podname string, servicename string) {

	apiserver.SaveRelationServicetoPod(servicename, podname)
	relations.ServicetoPodRela[servicename] = append(relations.ServicetoPodRela[servicename], servicename)

}
func RemovePodfromNode(podname string, nodename string) {
	for i, pod := range relations.NodetoPodRela[nodename] {
		if pod == podname {
			relations.NodetoPodRela[nodename] = append(
				relations.NodetoPodRela[nodename][:i],
				relations.NodetoPodRela[nodename][i+1:]...)
		}
	}
}
