package etcd

import (
	"fmt"
	mvccpb2 "go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func defaultHandler(event *clientv3.Event) error {
	switch event.Type {
	case mvccpb2.PUT:
		fmt.Println("修改为：", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
	case mvccpb2.DELETE:
		fmt.Println("删除了：", "Revision:", event.Kv.ModRevision, " key = ", string(event.Kv.Key))
	}
	return nil
}
