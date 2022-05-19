package constant

const (
	EtcdIPAddr string = "localhost:2379"

	ControllerPrefix string = "controller"
	RegistryPrefix   string = "registry"
	SchedulerPrefix  string = "Scheduler"

	PodSourceName     string = "pods"
	NodeSourceName    string = "nodes"
	ServiceSourceName string = "services"

	DefaultNameSpace string = "default"

	CREATE string = "create"
	DELETE string = "delete"
	UPDATE string = "update"
)
