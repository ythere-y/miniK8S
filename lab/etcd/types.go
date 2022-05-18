package etcd

type Options struct {
	Prefix     string
	SourceType string
	NameSpace  string
	Name       string
}
type Option func(*string)

type Handler func(key string, value string)

func SetPrefix(prefix string) Option {
	return func(this *string) {
		*this = "/" + (*this) + prefix
	}
}
func SetSourceType(sourceType string) Option {
	return func(this *string) {
		*this = (*this) + "/" + sourceType
	}
}

func SetNameSpace(nameSpace string) Option {
	return func(this *string) {
		*this = (*this) + "/" + nameSpace
	}
}
func SetName(name string) Option {
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
