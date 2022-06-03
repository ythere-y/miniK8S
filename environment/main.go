package environment

import (
	"fmt"
	"minik8s/config"
	"os/exec"
	"time"
)

func Main() {
	var (
		command *exec.Cmd
		err     error
	)
	//command = exec.Command("cd", "'E:\\'")
	//err = command.Run()
	//if err != nil {
	//	panic(err)
	//}

	command = exec.Command("notepad")
	err = command.Run()
	if err != nil {
		panic(err)
	}

	command = exec.Command("pwd")
	err = command.Run()
	if err != nil {
		panic(err)
	}

	tick := time.Tick(time.Second)
	for range tick {
		fmt.Println("...")
	}
}

func EtcdStartUp() error {
	var (
		err error
	)
	fmt.Printf("[Etcd] start up, IP = %v\n", config.Configs.EtcdIp)

	err = nil
	return err

}
func FlannelStartUp(etcdip string) error {
	var (
		err error
	)
	fmt.Printf("[Flannel] start up, etcdIP = %v\n", etcdip)
	err = nil
	return err
}
