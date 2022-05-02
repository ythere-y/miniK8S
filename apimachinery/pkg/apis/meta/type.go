package meta

// 代码源自：k8s.io/apimachinery/pkg/apis/meta/v1/types.go
// 来自apimachinery仓库
type TypeMeta struct {
	// Kind is a string value representing the REST resource this object represents.
	// Servers may infer this from the endpoint the client submits requests to.
	// Cannot be updated.
	// In CamelCase.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	// Kind 是一个字符串值，表示此对象所代表的 REST 资源。//服务器可以从客户端向其提交请求的端点推断出这一点。//无法更新。
	Kind string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`

	// APIVersion defines the versioned schema of this representation of an object.
	// Servers should convert recognized schemas to the latest internal value, and
	// may reject unrecognized values.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	// +optional
	// APIVersion定义对象表示的版本化架构。服务器应将已识别的架构转换为最新的内部值，并且可能会拒绝未识别的值。
	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`
}

// 代码源自k8s.io/apimachinery/pkg/apis/meta/v1/meta.go

//func (obj *TypeMeta) GetObjectKind() schema.ObjectKind { return obj }

// metav1.TypeMeta实现了schema.ObjectKind
// SetGroupVersionKind satisfies the ObjectKind interface for all objects that embed TypeMeta
// 对于所有嵌入 TypeMeta 的对象，SetGroupVersionKind 满足 ObjectKind 接口
//func (obj *TypeMeta) SetGroupVersionKind(gvk schema.GroupVersionKind) {
//	obj.APIVersion, obj.Kind = gvk.ToAPIVersionAndKind()
//}

// GroupVersionKind satisfies the ObjectKind interface for all objects that embed TypeMeta
// 对于所有嵌入 TypeMeta 的对象，GroupVersionKind 满足 ObjectKind 接口
//func (obj *TypeMeta) GroupVersionKind() schema.GroupVersionKind {
//	return schema.FromAPIVersionAndKind(obj.APIVersion, obj.Kind)
//}

// 代码源自k8s.io/apimachinery/pkg/runtime/schema/interfaces.go
// ObjectKind是接口，两个接口函数是GroupVersionKind类型的setter和getter
// 这个接口的作用下一章再介绍
type ObjectKind interface {
	//SetGroupVersionKind(kind GroupVersionKind)
	//GroupVersionKind() GroupVersionKind
}

// 代码源自：k8s.io/apimachinery/pkg/apis/meta/v1/types.go
// 同样来自apimachinery仓库，ObjectMeta是xxx.yaml中metadata字段，平时我们填写的metadata
// 一般只有name、label，其实ObjectMeta字段还是有很多内容了，让笔者逐一介绍一下。
type ObjectMeta struct {
	// 对象的名字应该不用介绍了。
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	// 如果Name为空，系统这为该对象生成一个唯一的名字。
	GenerateName string `json:"generateName,omitempty" protobuf:"bytes,2,opt,name=generateName"`
	// 命名空间，在平时学习、调试的时候很少用，但是在发布的时候会经常用。
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`
	// 对象的URL，由系统生成。
	SelfLink string `json:"selfLink,omitempty" protobuf:"bytes,4,opt,name=selfLink"`
	// 对象的唯一ID，由系统生成。
	//UID types.UID `json:"uid,omitempty" protobuf:"bytes,5,opt,name=uid,casttype=k8s.io/kubernetes/pkg/types.UID"`
	// 资源版本，这是一个非常有意思且变量，版本可以理解为对象在时间轴上的一个时间节点，代表着对象最后
	// 一次更新的时刻。如果说Name是在Namespace空间下唯一，那么ResourceVersion则是同名、同类型
	// 对象时间下唯一。因为同名对象在不同时间可能会更新、删除再添加，在比较两个对象谁比较新的情况
	// 非常有用，比如Watch。
	ResourceVersion string `json:"resourceVersion,omitempty" protobuf:"bytes,6,opt,name=resourceVersion"`
	// 笔者很少关注这个值，所以也不太了解是干什么的，读者感兴趣可以自己了解一下。
	Generation int64 `json:"generation,omitempty" protobuf:"varint,7,opt,name=generation"`
	// 对象创建时间，由系统生成。
	//CreationTimestamp Time `json:"creationTimestamp,omitempty" protobuf:"bytes,8,opt,name=creationTimestamp"`
	// 对象删除时间，指针类型说明是可选的，当指针不为空的时候说明对象被删除了，也是由系统生成
	//DeletionTimestamp *Time `json:"deletionTimestamp,omitempty" protobuf:"bytes,9,opt,name=deletionTimestamp"`
	// 对象被删除前允许优雅结束的时间，单位为秒。
	DeletionGracePeriodSeconds *int64 `json:"deletionGracePeriodSeconds,omitempty" protobuf:"varint,10,opt,name=deletionGracePeriodSeconds"`
	// 对象标签，这个是我们经常用的，不用多解释了
	Labels map[string]string `json:"labels,omitempty" protobuf:"bytes,11,rep,name=labels"`
	// 批注，这个和标签很像，但是用法不同，比如可以用来做配置。
	Annotations map[string]string `json:"annotations,omitempty" protobuf:"bytes,12,rep,name=annotations"`
	// 该对象依赖的对象类表，如果这些依赖对象全部被删除了，那么该对象也会被回收。如果该对象对象被
	// 某一controller管理，那么类表中有一条就是指向这个controller的。例如Deployment对象的
	// OwnerReferences有一条就是指向DeploymentController的。
	//OwnerReferences []OwnerReference `json:"ownerReferences,omitempty" patchStrategy:"merge" patchMergeKey:"uid" protobuf:"bytes,13,rep,name=ownerReferences"`
	// 下面这几个变量笔者没有了解过，等笔者知道了再来更新文章吧。
	// 表示删除前需要做什么清理操作，若DeletionTimestamp被设置时间，同时Finalizers不为空
	// 那么相应控制器将会检测到 Finalizers 中有自己设置的标签，就会做清理工作，做完清理掉自己设置的标签
	Finalizers  []string `json:"finalizers,omitempty" patchStrategy:"merge" protobuf:"bytes,14,rep,name=finalizers"`
	ClusterName string   `json:"clusterName,omitempty" protobuf:"bytes,15,opt,name=clusterName"`
}

// 代码源自k8s.io/apimachinery/pkg/apis/meta/v1/meta.go，GetObjectMeta是MetaAccessor
// 的接口函数，这个函数说明了ObjectMeta实现了MetaAccessor。
//func (obj *ObjectMeta) GetObjectMeta() Object { return obj }

// 下面所有的函数是接口Object的ObjectMeta实现，可以看出来基本就是setter/getter方法。纳尼？
// 又一个Object(这个Object是metav1.Object，本章节简写为Object)?Object是API对象公共属性
// (meta信息)的抽象，下面的函数是Object所有函数的实现，因为功能比较简单，笔者就不一一注释了。
//func (meta *ObjectMeta) GetNamespace() string                { return meta.Namespace }
//func (meta *ObjectMeta) SetNamespace(namespace string)       { meta.Namespace = namespace }
//func (meta *ObjectMeta) GetName() string                     { return meta.Name }
//func (meta *ObjectMeta) SetName(name string)                 { meta.Name = name }
//func (meta *ObjectMeta) GetGenerateName() string             { return meta.GenerateName }
//func (meta *ObjectMeta) SetGenerateName(generateName string) { meta.GenerateName = generateName }
//func (meta *ObjectMeta) GetUID() types.UID                   { return meta.UID }
//func (meta *ObjectMeta) SetUID(uid types.UID)                { meta.UID = uid }
//func (meta *ObjectMeta) GetResourceVersion() string          { return meta.ResourceVersion }
//func (meta *ObjectMeta) SetResourceVersion(version string)   { meta.ResourceVersion = version }
//func (meta *ObjectMeta) GetGeneration() int64                { return meta.Generation }
//func (meta *ObjectMeta) SetGeneration(generation int64)      { meta.Generation = generation }
//func (meta *ObjectMeta) GetSelfLink() string                 { return meta.SelfLink }
//func (meta *ObjectMeta) SetSelfLink(selfLink string)         { meta.SelfLink = selfLink }
//func (meta *ObjectMeta) GetCreationTimestamp() Time          { return meta.CreationTimestamp }
//func (meta *ObjectMeta) SetCreationTimestamp(creationTimestamp Time) {
//	meta.CreationTimestamp = creationTimestamp
//}
//func (meta *ObjectMeta) GetDeletionTimestamp() *Time { return meta.DeletionTimestamp }
//func (meta *ObjectMeta) SetDeletionTimestamp(deletionTimestamp *Time) {
//	meta.DeletionTimestamp = deletionTimestamp
//}
//func (meta *ObjectMeta) GetDeletionGracePeriodSeconds() *int64 {
//	return meta.DeletionGracePeriodSeconds
//}
//func (meta *ObjectMeta) SetDeletionGracePeriodSeconds(deletionGracePeriodSeconds *int64) {
//	meta.DeletionGracePeriodSeconds = deletionGracePeriodSeconds
//}
//func (meta *ObjectMeta) GetLabels() map[string]string                 { return meta.Labels }
//func (meta *ObjectMeta) SetLabels(labels map[string]string)           { meta.Labels = labels }
//func (meta *ObjectMeta) GetAnnotations() map[string]string            { return meta.Annotations }
//func (meta *ObjectMeta) SetAnnotations(annotations map[string]string) { meta.Annotations = annotations }
//func (meta *ObjectMeta) GetFinalizers() []string                      { return meta.Finalizers }
//func (meta *ObjectMeta) SetFinalizers(finalizers []string)            { meta.Finalizers = finalizers }
//func (meta *ObjectMeta) GetOwnerReferences() []OwnerReference         { return meta.OwnerReferences }
//func (meta *ObjectMeta) SetOwnerReferences(references []OwnerReference) {
//	meta.OwnerReferences = references
//}
func (meta *ObjectMeta) GetClusterName() string            { return meta.ClusterName }
func (meta *ObjectMeta) SetClusterName(clusterName string) { meta.ClusterName = clusterName }

//func (meta *ObjectMeta) GetManagedFields() []ManagedFieldsEntry { return meta.ManagedFields }
//func (meta *ObjectMeta) SetManagedFields(managedFields []ManagedFieldsEntry) {
//	meta.ManagedFields = managedFields
//}

// 代码源自：k8s.io/apimachinery/pkg/apis/meta/v1/types.go
type ListMeta struct {
	// 下面这两个变量在ObjectMeta相同，不多解释
	SelfLink        string `json:"selfLink,omitempty" protobuf:"bytes,1,opt,name=selfLink"`
	ResourceVersion string `json:"resourceVersion,omitempty" protobuf:"bytes,2,opt,name=resourceVersion"`
	// 在列举对象的时候可能会有非常多的对象，kubernetes支持分页获取，类似于SQL的limit，当对象总量
	// 多于单页的总量的时候，这个变量就会被设置。它用来告知用户需要继续获取，并且它包含了下次获取的
	// 起始位置。
	Continue string `json:"continue,omitempty" protobuf:"bytes,3,opt,name=continue"`
	// 从字面意思也能理解，就是还剩多少个对象，这个和Continue是配合使用的，当Continue被设置了
	// 这个变量就不会为空，用来告诉用户还有多少对象没有获取。
	RemainingItemCount *int64 `json:"remainingItemCount,omitempty" protobuf:"bytes,4,opt,name=remainingItemCount"`
}

// 代码源自k8s.io/apimachinery/pkg/apis/meta/v1/meta.go
// 下面所有的函数是接口ListInterface的ListMeta实现，ListInterface是API对象列表公共属性(meta)
// 的抽象。
//func (meta *ListMeta) GetResourceVersion() string        { return meta.ResourceVersion }
//func (meta *ListMeta) SetResourceVersion(version string) { meta.ResourceVersion = version }
//func (meta *ListMeta) GetSelfLink() string               { return meta.SelfLink }
//func (meta *ListMeta) SetSelfLink(selfLink string)       { meta.SelfLink = selfLink }
//func (meta *ListMeta) GetContinue() string               { return meta.Continue }
//func (meta *ListMeta) SetContinue(c string)              { meta.Continue = c }
//func (meta *ListMeta) GetRemainingItemCount() *int64     { return meta.RemainingItemCount }
//func (meta *ListMeta) SetRemainingItemCount(c *int64)    { meta.RemainingItemCount = c }

// k8s.io/apimachinery@v0.20.2/pkg/apis/meta/v1/meta.go
// Object lets you work with object metadata from any of the versioned or
// internal API objects. Attempting to set or retrieve a field on an object that does
// not support that field (Name, UID, Namespace on lists) will be a no-op and return
// a default value.
type Object interface {
	GetNamespace() string
	SetNamespace(namespace string)
	GetName() string
	SetName(name string)
	GetGenerateName() string
	SetGenerateName(name string)
	//GetUID() types.UID
	//SetUID(uid types.UID)
	GetResourceVersion() string
	SetResourceVersion(version string)
	GetGeneration() int64
	SetGeneration(generation int64)
	GetSelfLink() string
	SetSelfLink(selfLink string)
	//GetCreationTimestamp() Time
	//SetCreationTimestamp(timestamp Time)
	//GetDeletionTimestamp() *Time
	//SetDeletionTimestamp(timestamp *Time)
	GetDeletionGracePeriodSeconds() *int64
	SetDeletionGracePeriodSeconds(*int64)
	GetLabels() map[string]string
	SetLabels(labels map[string]string)
	GetAnnotations() map[string]string
	SetAnnotations(annotations map[string]string)
	GetFinalizers() []string
	SetFinalizers(finalizers []string)
	//GetOwnerReferences() []OwnerReference
	//SetOwnerReferences([]OwnerReference)
	GetClusterName() string
	SetClusterName(clusterName string)
	//GetManagedFields() []ManagedFieldsEntry
	//SetManagedFields(managedFields []ManagedFieldsEntry)
}

// ListInterface lets you work with list metadata from any of the versioned or
// internal API objects. Attempting to set or retrieve a field on an object that does
// not support that field will be a no-op and return a default value.
// TODO: move this, and TypeMeta and ListMeta, to a different package
type ListInterface interface {
	GetResourceVersion() string
	SetResourceVersion(version string)
	GetSelfLink() string
	SetSelfLink(selfLink string)
	GetContinue() string
	SetContinue(c string)
	GetRemainingItemCount() *int64
	SetRemainingItemCount(c *int64)
}

// ServiceSpec 描述了 Service 资源的属性
type ServiceSpec struct {
	// Service 暴露的所有端口
	//Ports []ServicePort `json:"ports,omitempty" patchStrategy:"merge" patchMergeKey:"port" protobuf:"bytes,1,rep,name=ports"`

	// 将 Service 的流量路由给拥有标签匹配 Selector 中定义的标签的 Pod。如果不指定或者为空，那么就假设这个
	// 服务有外部的进程来管理 Service 对应的 Endpoints，Kubernetes 在这种情况下不会去修改这个 Endpoints。
	// Selector 仅在 ServiceType 为 ClusterIP、NodePort 和 LoadBalancer 时生效。如果 ServiceType
	// 为 ExternalName，那么 Selector 值会被忽略。
	Selector map[string]string `json:"selector,omitempty" protobuf:"bytes,2,rep,name=selector"`

	// ClusterIP 是集群分配给 Service 的一个随机的 IP 地址。如果这个 IP 地址是手动指定的，那么要看这个 IP 是否已经
	// 被其它服务所占用。如果没有占用，那么就会分配这个 Service，否则 Service 的创建会失败。
	// 这个值不可以通过 Update 请求来更新。可选值为 None、空字符串和合法的 IP 地址。
	// 可以为无头 Service 指定 None，这个时候不需要转发流量。
	// ClusterIP 仅在 ServiceType 为 ClusterIP、NodePort 和 LoadBalancer 时生效。如果 ServiceType
	// 为 ExternalName，那么 ClusterIP 值会被忽略。
	//ClusterIP string `json:"clusterIP,omitempty" protobuf:"bytes,3,opt,name=clusterIP"`

	// Type 参数定义如何暴露这个 Service。 默认设置为 ClusterIP。可选值为 ExternalName、ClusterIP、NodePort 和
	// LoadBalancer。其中 ExternalName 表示使用指定的 externalName。
	// ClusterIP 表示为 Service 分配一个集群内部的 IP 并且将流量负载到 Service 对应的 Endpoints 去。如果 ClusterIP 设置为 None，
	// 将不会给 Service 分配一个虚拟 IP，Endpoints 就将作为一组端点暴露出来，而不是一个稳定的 IP。
	// NodePort 表示除了给 Service 分配一个集群内的 IP 之外，还会在每个 Node 上面为这个 Service 暴露一个映射端口。通过这个暴露出来的
	// 端口，Node 上的集群外的服务可以访问集群内的 Service。
	// LoadBalancer 基于 NodePort 之上，通过创建一个外部的负载均衡器，然后通过 NodePort 的端口路由到集群内部的 Service。
	//Type ServiceType `json:"type,omitempty" protobuf:"bytes,4,opt,name=type,casttype=ServiceType"`

	// 指定一组外部允许访问集群内服务的 IP。这些 IP 本身不属于 Kubernetes 系统。
	// 用户需要自己确保这些 IP 是允许访问集群内服务的。一个常见的例子是外部的负载均衡器。
	ExternalIPs []string `json:"externalIPs,omitempty" protobuf:"bytes,5,rep,name=externalIPs"`

	// 用来设置 Session 的亲和性。可选值为 ClientIP 和 None。默认为 None。
	//SessionAffinity ServiceAffinity `json:"sessionAffinity,omitempty" protobuf:"bytes,7,opt,name=sessionAffinity,casttype=ServiceAffinity"`

	// 只适用于 Service Type 为 LoadBalancer 的情况。用来指定负载均衡器的 IP 地址。
	LoadBalancerIP string `json:"loadBalancerIP,omitempty" protobuf:"bytes,8,opt,name=loadBalancerIP"`

	// 如果指定的参数被平台所支持，将控制只允许指定的 Client IP 的流量经过云厂商负载均衡器。
	LoadBalancerSourceRanges []string `json:"loadBalancerSourceRanges,omitempty" protobuf:"bytes,9,opt,name=loadBalancerSourceRanges"`

	// 该参数指定 kubedns 或者其他 DNS 返回的一个该 Service 的 CNAME 记录。这个中间不存在任何代理。
	ExternalName string `json:"externalName,omitempty" protobuf:"bytes,10,opt,name=externalName"`

	// 该参数表示 Service 期望如何路由外部的流量到 Node 或者是 Cluster 层面的 Endpoints。
	// 如果指定为 Local，那么会保留客户端的 IP 地址，对 LoadBalancer 和 NodePort 类型的服务来说，会避免第二跳，
	// 潜在的风险是导致流量负载的不均衡。
	// 如果指定为 Cluster，那么客户端的 IP 地址将被忽略，有可能导致流量负载的第二跳，但是可以保证整体流量的负载均衡。
	//ExternalTrafficPolicy ServiceExternalTrafficPolicyType `json:"externalTrafficPolicy,omitempty" protobuf:"bytes,11,opt,name=externalTrafficPolicy"`

	// 该参数指定 Service 健康检查的端口。如果没有指定的话，将自动创建一个关联到分配的 NodePort 的检查端口。
	// 只在 Service Type 为 LoadBalancer 和 ExternalTrafficPolicy 设置为 Local 的时候生效。
	HealthCheckNodePort int32 `json:"healthCheckNodePort,omitempty" protobuf:"bytes,12,opt,name=healthCheckNodePort"`

	// 当设置为 true 时，表示 Endpoints 中的 notReadyAddresses 对应的端点也需要发布出来。默认为 false。
	PublishNotReadyAddresses bool `json:"publishNotReadyAddresses,omitempty" protobuf:"varint,13,opt,name=publishNotReadyAddresses"`

	// 包含 Session 亲和性相关配置
	//SessionAffinityConfig *SessionAffinityConfig `json:"sessionAffinityConfig,omitempty" protobuf:"bytes,14,opt,name=sessionAffinityConfig"`
}

// ServicePort 定义了 Service 暴露的端口。
type ServicePort struct {
	// Service 内部的端口名称，ServiceSpec 内部的端口名称必须唯一。
	// 这个名称和 EndpointPort 对象名称一致。
	// 如果这个 Service 只暴露了一个端口，那么这个值是可选的。
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`

	// 该端口的 IP 协议。支持 TCP、UDP 和 SCTP。默认为 TCP 。
	//Protocol Protocol `json:"protocol,omitempty" protobuf:"bytes,2,opt,name=protocol,casttype=Protocol"`

	// Service 暴露的端口
	Port int32 `json:"port" protobuf:"varint,3,opt,name=port"`

	// Service 所指向的 Pod 的端口，这个端口可以是一个整型的端口号，也可以是端口的名称。
	// 端口的范围为 1 ～ 65535。 如果这个值是一个字符串，那么这个值会被当作目标 Port 的容器端口名称。
	// 如果没有指定，则使用和上面的 Port 一样的值。
	// 当 ClusterIP = None 的时候，这个值会被忽略。
	//TargetPort intstr.IntOrString `json:"targetPort,omitempty" protobuf:"bytes,4,opt,name=targetPort"`

	// 当 ServiceType 为 NodePort 或者 LoadBalancer 的时候，这个端口表示 Service 暴露在 Node 上面的端口。
	// 这个值通常是由系统分配的。如果是手动指定的话，当这个端口可用时会被分配给 Service，如果不可用那么 Service
	// 创建会失败。
	// 默认情况下，如果 Service 的 ServiceType 需要这样的一个端口的话，系统会自动分配一个。
	NodePort int32 `json:"nodePort,omitempty" protobuf:"varint,5,opt,name=nodePort"`
}

// ServiceStatus represents the current status of a service
type ServiceStatus struct {
	// LoadBalancer contains the current status of the load-balancer,
	// if one is present.
	// +optional
	//LoadBalancer LoadBalancerStatus

	// Current service condition
	// +optional
	//Conditions []metav1.Condition
}
