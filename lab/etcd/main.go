package etcd

import (
	"context"
	"flag"
	"fmt"
	mvccpb2 "go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"minik8s/constant"
	"minik8s/utils"
	"sync"
	"time"
)

// etcd client put/get demo
// use etcd/clientv3

var (
	syncCount = 0
	addr      = flag.String("addr", "http://127.0.0.1:2379", "etcd address")
	queueName = flag.String("name", "my-test-queue", "queue name")
	mtx       sync.Mutex
)

func defaultPutHander(key string, value string) {

}
func LabMain() {
	EasyPutGetTest()
	//MessagingTest()
	//WatchTest()
	//
	//go TestRemoteIp(ip, "hello", "world")
	//go Watch("minik")
	//go Watch("t2")
	//go Watch("t3")
	//go SyncPutTest("minik", " --help")
	//go SyncPutTest("minik", " pod -h")
	//go SyncPutTest("minik", "pod -h")
	//go SyncPutTest("t2", "pod -h")

	utils.HoldPro()
}

func TestRemoteIp(ip string, key string, input string) {
	time.Sleep(time.Second * 3)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{constant.EtcdIPAddr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")

	defer cli.Close()
	// put
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, key, input)
	//cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
}

func GetWithPrefix(key string) (*clientv3.GetResponse, error) {
	var (
		err    error
		cli    *clientv3.Client
		getRsp *clientv3.GetResponse
		ctx    context.Context
		cancel context.CancelFunc
	)
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{constant.EtcdIPAddr},
		DialTimeout: 5 * time.Second,
	})

	utils.HandleError("connect to etcd failed", err)
	defer cli.Close()
	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	getRsp, err = cli.Get(ctx, key, clientv3.WithPrefix())
	cancel()

	return getRsp, err
}
func GetNormal(key string) (*clientv3.GetResponse, error) {
	var (
		err    error
		cli    *clientv3.Client
		getRsp *clientv3.GetResponse
		ctx    context.Context
		cancel context.CancelFunc
	)
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{constant.EtcdIPAddr},
		DialTimeout: 5 * time.Second,
	})

	utils.HandleError("connect to etcd failed", err)
	defer cli.Close()
	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	getRsp, err = cli.Get(ctx, key)
	cancel()

	return getRsp, err
}

func Get(key string) ([]string, error) {
	var (
		//kv     clientv3.KV
		res    []string
		err    error
		cli    *clientv3.Client
		getRsp *clientv3.GetResponse
		ctx    context.Context
		cancel context.CancelFunc
	)
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{constant.EtcdIPAddr},
		DialTimeout: 5 * time.Second,
	})

	utils.HandleError("connect to etcd failed", err)

	fmt.Println("connect to etcd success")
	/* 另一种方式
	kv = clientv3.NewKV(cli)

	getRsp, err = kv.Get(context.TODO(), "/demo/A", clientv3.WithPrefix())

	utils.HandleError("get error", err)

	for _, resp := range getRsp.Kvs {
		fmt.Printf("key: %s, value:%s\n", string(resp.Key), string(resp.Value))
	}
	*/

	defer cli.Close()
	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	getRsp, err = cli.Get(ctx, key)
	cancel()

	utils.HandleError("get operation error", err)
	// 遍历得到的所有结果
	for _, ev := range getRsp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
		res = append(res, string(ev.Value))
	}

	return res, err
}

func EasyPutGetTest() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{constant.EtcdIPAddr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()
	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, "/test/a", "a-hello world")
	_, err = cli.Put(ctx, "/test/b", "b-hello world")
	_, err = cli.Put(ctx, "/test", "no-hello world")
	_, err = cli.Put(ctx, "/test/a/in_a", "in_a -hello world")
	_, err = cli.Put(ctx, "/test/a/in a", "in block a -hello world")
	_, err = cli.Put(ctx, "test/a", "head no hello")
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}

	// get
	fmt.Println("test 1 ,prefix key = /test")
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "/test", clientv3.WithPrefix())
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("[key] %s :[value] %s\n", ev.Key, ev.Value)
	}

	fmt.Println("test 2, prefix key = /test/")
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err = cli.Get(ctx, "/test/", clientv3.WithPrefix())
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("[key] %s :[value] %s\n", ev.Key, ev.Value)
	}
}

//etcd 实现分布式队列

func SyncPut(key string, value string) {
	go Put(key, value)
}
func Put(key string, value string) {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{constant.EtcdIPAddr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")

	defer cli.Close()
	// put
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, key, value)
	//cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
	fmt.Printf("Put op : key = %v, val = %v\n", key, value)
}

func SyncPutTest(key string, input string) {
	time.Sleep(time.Second * 3)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{constant.EtcdIPAddr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")

	defer cli.Close()
	// put
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, key, input)
	//cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
}

func SyncWatch(name string, handler Handler) {
	go WatchWithFunc(name, handler)
}

func WatchWithFunc(name string, handler Handler) {
	config := clientv3.Config{
		Endpoints:   []string{constant.EtcdIPAddr},
		DialTimeout: 5 * time.Second,
	}
	cli, err := clientv3.New(config)
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()

	watchRespChan := cli.Watch(context.Background(), name, clientv3.WithPrefix()) // <-chan WatchResponse
	for watchResp := range watchRespChan {
		mtx.Lock()
		for _, event := range watchResp.Events {
			fmt.Printf("Type: %s\t Key:%s \n", event.Type, event.Kv.Key)
			switch event.Type {
			case mvccpb2.PUT:
				fmt.Println("修改为：", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb2.DELETE:
				fmt.Println("删除了：", "Revision:", event.Kv.ModRevision)
			}
			handler(event)
		}
		mtx.Unlock()
	}
}

func Watch(name string, handler Handler) {
	config := clientv3.Config{
		Endpoints:   []string{constant.EtcdIPAddr},
		DialTimeout: 5 * time.Second,
	}
	cli, err := clientv3.New(config)
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()

	watchRespChan := cli.Watch(context.Background(), name, clientv3.WithPrefix()) // <-chan WatchResponse
	for watchResp := range watchRespChan {
		mtx.Lock()
		for _, event := range watchResp.Events {
			fmt.Printf("Type: %s\t Key:%s \n", event.Type, event.Kv.Key)
			switch event.Type {
			case mvccpb2.PUT:
				fmt.Println("修改为：", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb2.DELETE:
				fmt.Println("删除了：", "Revision:", event.Kv.ModRevision)
			}
			syncCount++

		}
		mtx.Unlock()
	}
}

func ServiceWatchSync() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{constant.EtcdIPAddr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()
	// watch key:q1mi change
	rch := cli.Watch(context.Background(), "service") // <-chan WatchResponse
	for wresp := range rch {
		mtx.Lock()
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			fmt.Printf("syncCount = %v, and start sleeping\n", syncCount)
			time.Sleep(time.Second * 10)
		}
		mtx.Unlock()
	}
}
