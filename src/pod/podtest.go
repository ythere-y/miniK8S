package pod

func main() {
	file := "yamltest1.yaml"
	CreatePod(file)
	GetAllPodInfo()
	RunPodByName("testpod")
}
