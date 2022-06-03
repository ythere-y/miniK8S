package etcd

import (
	"context"
	"fmt"
	mvccpb2 "go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"minik8s/config"
	"minik8s/constant"
	"minik8s/utils"
	"sync"
	"time"
)

// etcd client put/get demo
// use etcd/clientv3

var (
	syncCount = 0
	mtx       sync.Mutex
)

func LabMain() {
	//EasyPutGetTest()
	DeleteWatchTest()
	utils.HoldPro()

}

func GetWithPrefix(key string) (*clientv3.GetResponse, error) {
	var (
		err    error
		cli    *clientv3.Client
		getRsp *clientv3.GetResponse
		ctx    context.Context
		cancel context.CancelFunc
	)
	cli = connectEtcd()
	defer cli.Close()
	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	getRsp, err = cli.Get(ctx, key, clientv3.WithPrefix())
	cancel()

	return getRsp, err
}
func GetValueWithPrefix(key string) (string, error) {
	var (
		err    error
		cli    *clientv3.Client
		getRsp *clientv3.GetResponse
		ctx    context.Context
		cancel context.CancelFunc
	)
	cli = connectEtcd()
	defer cli.Close()
	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	getRsp, err = cli.Get(ctx, key, clientv3.WithPrefix())
	cancel()
	if getRsp.Count == 0 {
		return "", err
	}
	return string(getRsp.Kvs[0].Value), err
}

func GetValue(key string) (string, error) {
	var (
		err    error
		cli    *clientv3.Client
		getRsp *clientv3.GetResponse
		ctx    context.Context
		cancel context.CancelFunc
	)
	cli = connectEtcd()
	defer cli.Close()
	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	getRsp, err = cli.Get(ctx, key)
	cancel()
	if getRsp.Count == 0 {
		return "", err
	}
	return string(getRsp.Kvs[0].Value), err
}

func GetNormal(key string) (*clientv3.GetResponse, error) {
	var (
		err    error
		cli    *clientv3.Client
		getRsp *clientv3.GetResponse
		ctx    context.Context
		cancel context.CancelFunc
	)
	cli = connectEtcd()
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
	cli = connectEtcd()

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

//DeleteWatchTest
/*
测试函数
*/
func DeleteWatchTest() {
	go WatchWithFunc("/test/", defaultHandler)
	time.Sleep(2 * time.Second)
	Put("/test/h_1", "1__")
	Put("/test/h_2", "2__")
	Put("/test/h_3", "3__")
	Put("/test/h_4", "4")

	Delete("/test/h_1")
	utils.HoldPro()
}

//EasyPutGetTest
/*
测试函数
*/
func EasyPutGetTest() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{constant.EtcdIPAddr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
		return
	}
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

func connectEtcd() *clientv3.Client {
	if Endpoints == nil {
		Endpoints = []string{config.Configs.EtcdIp}
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   Endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return nil
	}
	fmt.Println("connect to etcd success")
	return cli
}

//Put
/*
给etcd中加入一个键值对
*/
func Put(key string, value string) {
	var (
		cli    *clientv3.Client
		err    error
		cancel context.CancelFunc
		ctx    context.Context
	)
	cli = connectEtcd()
	defer func(cli *clientv3.Client) {
		err := cli.Close()
		if err != nil {

		}
	}(cli)

	// put
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, key, value)
	cancel()
	//cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
	fmt.Printf("Put operation :\n[key] = %v \n[val] = %v\n", key, value)
}

//PutList
/*
给etcd中加入一个键值对
*/
func PutList(key []string, value []string) {
	var (
		cli     *clientv3.Client
		err     error
		cancel  context.CancelFunc
		ctx     context.Context
		listLen = len(key)
	)
	cli = connectEtcd()
	defer cli.Close()

	// put
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	for i := 0; i < listLen; i++ {
		_, err = cli.Put(ctx, key[i], value[i])
		//cancel()
		if err != nil {
			fmt.Printf("put to etcd failed, err:%v\n", err)
			cancel()
			return
		}
		fmt.Printf("Put operation : key = %v, val = %v\n", key, value)
	}
	cancel()
}

//Delete
/*
在etcd中删除一个key
*/
func Delete(key string) {
	var (
		cli    *clientv3.Client
		err    error
		cancel context.CancelFunc
		ctx    context.Context
	)
	cli = connectEtcd()
	defer func(cli *clientv3.Client) {
		err := cli.Close()
		if err != nil {

		}
	}(cli)

	// put
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Delete(ctx, key)
	cancel()
	//cancel()
	if err != nil {
		fmt.Printf("delete [key = %v] in etcd failed, err:%v\n", key, err)
		return
	}
	fmt.Printf("Delete operation : [key = %v] \n", key)
}

//DeleteListWithPrefix
/*
在etcd中删除一连串的key
*/
func DeleteListWithPrefix(key []string) {
	var (
		cli    *clientv3.Client
		err    error
		cancel context.CancelFunc
		ctx    context.Context
	)
	cli = connectEtcd()
	defer func(cli *clientv3.Client) {
		err := cli.Close()
		if err != nil {

		}
	}(cli)

	// del
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	for _, innerKey := range key {
		_, err = cli.Delete(ctx, innerKey, clientv3.WithPrefix())
		if err != nil {
			fmt.Printf("delete [key = %v] in etcd failed, err:%v\n", innerKey, err)
			cancel()
			return
		}
		fmt.Printf("Delete operation : [key = %v] \n", innerKey)
	}
	cancel()
}

//DeleteList
/*
在etcd中删除一连串的key
*/
func DeleteList(key []string) {
	var (
		cli    *clientv3.Client
		err    error
		cancel context.CancelFunc
		ctx    context.Context
	)
	cli = connectEtcd()
	defer func(cli *clientv3.Client) {
		err := cli.Close()
		if err != nil {

		}
	}(cli)

	// del
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	for _, innerKey := range key {
		_, err = cli.Delete(ctx, innerKey)
		if err != nil {
			fmt.Printf("delete [key = %v] in etcd failed, err:%v\n", innerKey, err)
			cancel()
			return
		}
		fmt.Printf("Delete operation : [key = %v] \n", innerKey)
	}
	cancel()
}

//DeleteWithPrefix
/*
在etcd中删除一个key,使用前缀搜索
*/
func DeleteWithPrefix(key string) {
	var (
		cli    *clientv3.Client
		err    error
		cancel context.CancelFunc
		ctx    context.Context
	)
	cli = connectEtcd()
	defer func(cli *clientv3.Client) {
		err := cli.Close()
		if err != nil {

		}
	}(cli)

	// put
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Delete(ctx, key, clientv3.WithPrefix())
	cancel()
	//cancel()
	if err != nil {
		fmt.Printf("delete [key = %v] in etcd failed, err:%v\n", key, err)
		return
	}
	fmt.Printf("Delete operation : [key = %v] \n", key)
}

//WatchWithFuncWithTime
/*
watch某个name，并且用handler进行后续处理,且超时之后自动放弃watch
*/
func WatchWithFuncWithTime(name string, handler Handler, fail FailOut, timeset time.Duration) {

	var (
		cli         *clientv3.Client
		err         error
		successFlag = false
	)
	cli = connectEtcd()
	defer cli.Close()

	ctx, cancle := context.WithCancel(context.Background())
	time.AfterFunc(timeset, func() {
		if successFlag == false {
			fail()

		}
		cancle()
		return
	})
	watchRespChan := cli.Watch(ctx, name, clientv3.WithPrefix()) // <-chan WatchResponse
	for watchResp := range watchRespChan {
		for _, event := range watchResp.Events {
			//err = defaultHandler(event)
			err = handler(event)
			if err != nil {
				return
			}
			successFlag = true
			cancle()
			return
		}
		//cancle()
		//return
	}

	fmt.Printf("watch on %v ended~!\n", name)
}

//WatchWithFunc
/*
watch某个name，并且用handler进行后续处理
*/
func WatchWithFunc(name string, handler Handler) {

	var (
		cli *clientv3.Client
		err error
	)
	cli = connectEtcd()
	defer cli.Close()

	watchRespChan := cli.Watch(context.Background(), name, clientv3.WithPrefix()) // <-chan WatchResponse
	for watchResp := range watchRespChan {
		for _, event := range watchResp.Events {
			//err = defaultHandler(event)
			err = handler(event)
			if err != nil {
				return
			}
		}
	}
}

func WatchWithCounter(name string) {
	var (
		cli *clientv3.Client
	)
	cli = connectEtcd()
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
