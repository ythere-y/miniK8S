package main

import (
	podm "minik8s/src/pod"
)

func main() {
	file := "podtest.yaml"
	podm.CreatePod(file)
	podm.GetAllPodInfo()
	podm.RunPodByName("testpod")
	podm.GetAllPodInfo()
}
