package environment

import (
	"fmt"
	"minik8s/config"
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
