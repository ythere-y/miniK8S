package service

// Container meta data in yaml

type ServiceYaml struct {
	Kind     string `yaml:"kind"`
	MetaData struct {
		Name string `yaml:"name"`
	}
	Selector []string `yaml:"selector,flow"`

	Port       int32 `yaml:"port"`
	TargetPort int32 `yaml:"targetPort"`
}

type MiniService struct {
	Kind     string
	MetaData struct {
		Name string
	}
	Selector map[string]string

	Port       int32
	TargetPort int32
}

func (mini MiniService) Build(yaml ServiceYaml) {
	mini.Kind = yaml.Kind
	mini.MetaData.Name = yaml.MetaData.Name
	mini.Port = yaml.Port
	mini.TargetPort = yaml.TargetPort
	for _, str := range yaml.Selector {
		//TODO:这部分的转换需要考虑一下
		mini.Selector[str] = "hello_but_nil"
	}
}

func (mini MiniService) Display() {
	println(mini)
}

type IntOrString struct {
	Type   Type   `protobuf:"varint,1,opt,name=type,casttype=Type"`
	IntVal int32  `protobuf:"varint,2,opt,name=intVal"`
	StrVal string `protobuf:"bytes,3,opt,name=strVal"`
}

// Type represents the stored type of IntOrString.
type Type int64

const (
	Int    Type = iota // The IntOrString holds an int.
	String             // The IntOrString holds a string.
)

func SerReadTest() {
	file := "servicetest.yaml"

	BuildService(file)
	GetAllServicInfo()
	CreateService(file)
}
