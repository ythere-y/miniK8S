package constant

const (
	EtcdIPAddr string = "localhost:2379"

	ControllerPrefix string = "controller"
	RegistryPrefix   string = "registry"
	SchedulerPrefix  string = "scheduler"
	RelationPrefix   string = "relation"

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

	EtcdSh        string = "./Scripts/etcdSh"
	EnvironmentSh string = "./Scripts/envStartUp.sh"

	MasterNodeFile string = "./config/masterNode.yaml"
)
