package controllerManager

type Relation struct {
	PodstoNodeRela   map[string]string
	NodetoPodRela    map[string][]string
	ServicetoPodRela map[string][]string
}
