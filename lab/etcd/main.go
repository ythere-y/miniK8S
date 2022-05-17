package etcd

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
	"log"
	"minik8s/cmd"
	"os"
	"strings"
	"sync"
	"time"
)

// etcd client put/get demo
// use etcd/clientv3

var syncCount = 0
var ip = "localhost:2379"

func LabMain() {
	EasyPutGetTest()
	//MessagingTest()
	//WatchTest()
	//
	go TestRemoteIp(ip, "hello", "world")
	//go SyncWatch("minik")
	//go SyncWatch("t2")
	//go SyncWatch("t3")
	//go SyncPutTest("minik --help")
	//go SyncPutTest("minik", " --help")
	//go SyncPutTest("minik", " pod -h")
	//go SyncPutTest("minik", "pod -h")

	holdPro()
}
func holdPro() {
	for true {
	}
}
func TestRemoteIp(ip string, key string, input string) {
	time.Sleep(time.Second * 3)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{ip},
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
func EasyPutGetTest() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	/*
		kv := clientv3.NewKV(cli)
		fmt.Println("Put part")
		kv.Put(context.TODO(), "demo/A/B1", "BB", clientv3.WithPrevKV())
		kv.Put(context.TODO(), "demo/A/B2", "CCC", clientv3.WithPrevKV())
		kv.Put(context.TODO(), "demo/A/B3", "DDDD", clientv3.WithPrevKV())

		fmt.Println("Get part")

		getRsp, err := kv.Get(context.TODO(), "/demo/A", clientv3.WithPrefix())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(getRsp.Kvs, getRsp.Count)
		for _, resp := range getRsp.Kvs {
			fmt.Printf("key: %s, value:%s\n", string(resp.Key), string(resp.Value))
		}
	*/

	defer cli.Close()
	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, "hh", "helloworld")
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "hh")
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}
}

//etcd 实现分布式队列

var (
	addr      = flag.String("addr", "http://127.0.0.1:2379", "etcd address")
	queueName = flag.String("name", "my-test-queue", "queue name")
	mtx       sync.Mutex
)

func MessagingTest() {
	flag.Parse()

	endpoints := strings.Split(*addr, ",")

	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	//创建获取队列
	q := recipe.NewQueue(cli, *queueName)

	//从命令行读取命令
	consol := bufio.NewScanner(os.Stdin)
	for consol.Scan() {
		action := consol.Text()
		items := strings.Split(action, " ")
		switch items[0] {
		case "push":
			if len(items) != 2 {
				fmt.Println("must set value to push")
				continue
			}
			q.Enqueue(items[1]) //入队
		case "pop":
			v, err := q.Dequeue() //出队
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(v) //输出出队元素
		case "quit", "exit":
			return
		default:
			fmt.Println("unknow action")
		}
	}

}
func SyncPutTest(key string, input string) {
	time.Sleep(time.Second * 3)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
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
func SyncWatch(name string) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()
	// watch key:q1mi change
	rch := cli.Watch(context.Background(), name) // <-chan WatchResponse
	for wresp := range rch {
		mtx.Lock()
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			syncCount++

			if name == "minik" {
				var setArgs []string
				setArgs = append(setArgs, string(ev.Kv.Value))
				cmd.RootCmd.SetArgs(strings.Fields(string((ev.Kv.Value))))

				err := cmd.RootCmd.Execute()
				if err != nil {
					fmt.Printf(err.Error())
				}
			}
			fmt.Printf("syncCount = %v, and start sleeping\n", syncCount)
			time.Sleep(time.Second * 2)
		}
		mtx.Unlock()
	}
}
func ServiceWatchSync() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
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
