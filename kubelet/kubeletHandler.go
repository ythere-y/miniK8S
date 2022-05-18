package kubelet

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func kubelethandler(event *clientv3.Event) {
	fmt.Printf("kubelet handling key = %v, value = %v\n", string(event.Kv.Key), string(event.Kv.Value))

}
