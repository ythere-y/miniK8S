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

	CREATE string = "create"
	DELETE string = "delete"
	UPDATE string = "update"
	STOP   string = "stop"

	StopFlag   string = "__stop__"
	RemoveFlag string = "__remove__"
	UpdateFlag string = "__update__"
)
