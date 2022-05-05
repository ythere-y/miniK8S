package pod

func PodTest() {
	file := "yamltest1.yaml"
	CreatePod(file)
	RunPodByName("testpod")
}

func main() {
	file := "yamltest1.yaml"
	CreatePod(file)
	GetAllPodInfo()
	//RunPodByName("testpod")
}

func PutPods() {
	file := "podtest.yaml"
	CreatePod(file)
	//GetAllPodInfo()

}
