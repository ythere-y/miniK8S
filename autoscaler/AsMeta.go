package autoscaler

type AutoScalerMetrics struct {
	Cpu float32
	Mem float32
}

type AutoScalerSpec struct {
	MinReplicas uint32
	MaxReplicas uint32
	Metrics     AutoScalerMetrics
}

type AutoScaler struct {
	Kind     string
	Name     string
	Workload string
	Spec     AutoScalerSpec
}

type AutoScalerYaml struct {
	Kind     string `yaml:"kind"`
	Name     string `yaml:"name"`
	Workload string `yaml:"workload"`
	Spec     struct {
		MinReplicas uint32 `yaml:"minreplicas"`
		MaxReplicas uint32 `yaml:"maxreplicas"`
		Metrics     struct {
			Cpu float32 `yaml:"cpu"`
			Mem float32 `yaml:"memory"`
		}
	}
}
