package etcd

import clientv3 "go.etcd.io/etcd/client/v3"

type Options struct {
	Prefix     string
	SourceType string
	NameSpace  string
	Name       string
}
type Option func(*string)

type FailOut func()

type Handler func(event *clientv3.Event) error

func JustAppend(context string) Option {
	return AppendName(context)
}
func SetPrefix(prefix string) Option {
	return AppendName(prefix)
}
func SetSourceType(sourceType string) Option {
	return AppendName(sourceType)
}

func SetNameSpace(nameSpace string) Option {
	return AppendName(nameSpace)
}
func SetName(name string) Option {
	return AppendName(name)
}
func SetNodeName(nodeName string) Option {
	return AppendName(nodeName)
}
func SetPodName(podName string) Option {
	return AppendName(podName)
}

func AppendName(name string) Option {
	return func(this *string) {
		*this = (*this) + "/" + name
	}
}

func SetKey(opts ...Option) string {
	key := ""

	for _, opt := range opts {
		opt(&key)
	}
	return key

}
