package main

import (
	"fmt"
	"github.com/bitfield/script"
	"minik8s/constant"
	"os"
	"os/exec"
)

func check(name string) string {
	goExecPath, err := exec.LookPath(name)

	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Go Executable: ", goExecPath)
	}
	return goExecPath

}

func startupEtcd(thisIP string) {
	var (
		//advertiseFlag string::
		frontArg     string
		advertiseArg string
		listenArg    string
		totalString  string
		cmdLine      string
		//count        int64
		get string
		err error
	)

	check("echo")
	totalString = "#!/bin/bash\n"

	frontArg = "etcd -name etcd-hc -data-dir /var/lib/etcd "
	advertiseArg = "http://" + thisIP + ":2379,http://127.0.0.1:2379 "
	listenArg += "--listen-client-urls http://" + thisIP + ":2379,http://127.0.0.1:2379 "
	totalString += frontArg + advertiseArg + listenArg + "\n"
	//totalString = "echo 'hello world~!'"
	if err != nil {
		panic(err)
	}
	// 检查etcd运行状态
	cmdLine = "etcdctl  --endpoints http://" + thisIP + ":2379 member list\n"
	totalString += cmdLine + "\n"
	// 设置flannel参数
	cmdLine = "etcdctl set  /coreos.com/network/config '{\"Network\": \"10.0.0.0/16\", \"SubnetLen\": 24, \"SubnetMin\": \"10.0.10.0\",\"SubnetMax\": \"10.0.20.0\", \"Backend\": {\"Type\": \"vxlan\"}}'"
	totalString += cmdLine + "\n"

	// 写入脚本
	fmt.Printf("ready to write [file = %v], context ->:\n%v\n", constant.EtcdSh, totalString)
	script.Echo(totalString).WriteFile(constant.EtcdSh)

	// 检验写入结果
	get, err = script.File(constant.EtcdSh).String()

	fmt.Printf("check the file :\n %v", get)

	return
	// 以下是执行该脚本内容的部分
	check("bash")

	goExecutable, _ := exec.LookPath("bash")

	cmdGoVer := &exec.Cmd{
		Path:   goExecutable,
		Args:   []string{goExecutable, constant.EtcdSh},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	fmt.Println(cmdGoVer.String())

	if err := cmdGoVer.Run(); err != nil {
		fmt.Println("Error: ", err)
	}

}

func main() {
	fmt.Print("hello world\n")
	startupEtcd("192.168.1.4")
}
