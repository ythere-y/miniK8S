package apiserver

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func apiserverHandler(event *clientv3.Event) error {
	var err error
	fmt.Printf("apiserver handling key = %v, value = %v\n", string(event.Kv.Key), string(event.Kv.Value))
	return err
}
