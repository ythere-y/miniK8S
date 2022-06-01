package node

type Addresses struct {
	HostName   string
	ExternalIp string
	InternalIp string
}

type Condition struct {
	OutOfDisk      bool // 如果节点上没有足够的可用空间来添加新的pod；否则为：False
	Ready          bool // 如果节点是健康的并准备好接收pod；False：如果节点不健康并且不接受pod；Unknown：如果节点控制器在过去40秒内没有收到node的状态报告。
	MemoryPressure bool // True：如果节点存储器上内存过低; 否则为：False。
	DiskPressure   bool // True：如果磁盘容量存在压力 - 也就是说磁盘容量低；否则为：False。
}

// OSName is the set of OS'es that can be used in OS.
type OSName string

// These are valid values for OSName
const (
	Linux   OSName = "linux"
	Windows OSName = "windows"
)

type Info struct {
	coreVersion     string
	kubernetVersion string
	dockerVersion   string
	OSName
}

type NodeStatus struct {
	Addresses
	Condition
	Capacity uint32
	Name     string
}
