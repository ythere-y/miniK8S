package constant

import "time"

const (
	EtcdIPAddr string = "localhost:2379"

	ControllerPrefix string = "controller"
	RegistryPrefix   string = "registry"
	SchedulerPrefix  string = "scheduler"
	RelationPrefix   string = "relation"
	ReplayPrefix     string = "reply"

	ReplayOK    string = "__OK__"
	ReplayERROR string = "__ERROR__"

	PodSourceName     string = "pods"
	NodeSourceName    string = "nodes"
	ServiceSourceName string = "services"
	ReplicaSourceName string = "replicaset"

	CREATE string = "create"
	DELETE string = "delete"
	UPDATE string = "update"
	STOP   string = "stop"

	StopFlag   string = "__stop__"
	RemoveFlag string = "__remove__"
	UpdateFlag string = "__update__"

	EtcdSh        string = "./Scripts/etcdSh.sh"
	EnvironmentSh string = "./Scripts/envStartUp.sh"
	TmpSh         string = "./Scripts/tmp.sh"
	FlannelSh     string = "./Scripts/flannelSh.sh"

	MasterNodeFile string = "./config/masterNode.yaml"
)

const (
	WaitReplyTime time.Duration = 10 * time.Second
)
