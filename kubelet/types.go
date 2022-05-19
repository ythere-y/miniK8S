package kubelet

import "minik8s/src/pod"

const (
	POD_PENDING   pod.PodStatusType = "PENDING"
	POD_RUNNING   pod.PodStatusType = "RUNNING"
	POD_SUCCEEDED pod.PodStatusType = "SUCCEEDED"
	POD_FAILED    pod.PodStatusType = "FAILED"
	POD_UNKNOWN   pod.PodStatusType = "UNKNOWN"
)
