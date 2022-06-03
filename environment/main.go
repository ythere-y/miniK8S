package environment

import (
	"fmt"
	"github.com/bitfield/script"
	"minik8s/config"
	"minik8s/constant"
	"os"
	"os/exec"
)

func Main() {
}

// TODO: 和startupetcd联合
func EtcdStartUp() error {
	var (
		err error
	)
	fmt.Printf("[Etcd] start up, IP = %v\n", config.Configs.EtcdIp)

	err = nil
	return err

}

// TODO: 和startupflannel联合
func FlannelStartUp(etcdip string) error {
	var (
		err error
	)
	fmt.Printf("[Flannel] start up, etcdIP = %v\n", etcdip)
	err = nil
	return err
}

// TODO: 和iptablesSet联合
func SetIptables(from string, to string) error {
	var (
		err error
	)
	fmt.Printf("[IPtables] set ~!\n")
	err = nil
	return err
}

func startupFlannel(targetIP string) {
	var (
		totalString string
		cmdLine     string
		get         string
		err         error
	)
	// 启动flannel
	cmdLine = "flannel -etcd-endpoints \"http://" + targetIP + ":4001,http://" + targetIP + "192.168.1.13:2379\""
	totalString += cmdLine + "\n"

	// 写入脚本
	script.Echo(totalString).WriteFile(constant.EnvironmentSh)

	// 检验写入结果
	get, err = script.File(constant.EnvironmentSh).String()

	if err != nil {
		panic(err)
	}

	fmt.Printf("check the file :\n %v", get)

	// 以下是执行该脚本内容的部分
	check("bash")

	goExecutable, _ := exec.LookPath("bash")

	cmdGoVer := &exec.Cmd{
		Path:   goExecutable,
		Args:   []string{goExecutable, constant.EnvironmentSh},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	fmt.Println(cmdGoVer.String())

	if err := cmdGoVer.Run(); err != nil {
		fmt.Println("Error: ", err)
	}

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
	script.Echo(totalString).WriteFile(constant.EtcdSh)

	// 检验写入结果
	get, err = script.File(constant.EtcdSh).String()

	fmt.Printf("check the file :\n %v", get)

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

func iptableSet() {

}

func check(name string) string {
	goExecPath, err := exec.LookPath(name)

	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Go Executable: ", goExecPath)
	}
	return goExecPath

}
