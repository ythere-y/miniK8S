package environment

import (
	"bytes"
	"fmt"
	"github.com/bitfield/script"
	"minik8s/config"
	"minik8s/constant"
	"os"
	"os/exec"
	"time"
)

func Main() {
}

func EtcdStartUp(thisIP string) error {
	var (
		err error
	)
	fmt.Printf("[Etcd] start up, IP = %v\n", config.Configs.EtcdIp)
	startupEtcd(thisIP)
	time.Sleep(10 * time.Second)
	setFlannelConfig(thisIP)
	err = nil
	return err

}

func FlannelStartUp(etcdip string) error {
	var (
		err error
	)
	fmt.Printf("[Flannel] start up, etcdIP = %v\n", etcdip)

	startupFlannel(etcdip)

	err = nil
	return err
}

func SetIptables(srcIP string, srcPort string, desIP string, desPort string) error {
	var (
		err error
	)
	fmt.Printf("[IPtables] set ~!\n")
	iptablesSet(srcIP, srcPort, desIP, desPort)
	err = nil
	return err
}

func RemoveIptables(srcIP string, srcPort string, desIP string, desPort string) error {
	var (
		err error
	)
	fmt.Printf("[IPtables] set ~!\n")
	iptablesSet(srcIP, srcPort, desIP, desPort)
	err = nil
	return err
}
func startupEtcd(thisIP string) {
	var (
		//advertiseFlag string::
		frontArg     string
		advertiseArg string
		listenArg    string
		totalString  string
		//count        int64
		err error
	)

	check("echo")
	totalString = "#!/bin/bash\n"

	frontArg = "etcd -name etcd-hc -data-dir /var/lib/etcd "
	advertiseArg = "--advertise-client-urls http://" + thisIP + ":2379,http://127.0.0.1:2379 "
	listenArg += "--listen-client-urls http://" + thisIP + ":2379,http://127.0.0.1:2379 "
	totalString += frontArg + advertiseArg + listenArg + " &\n"
	//totalString = "echo 'hello world~!'"
	if err != nil {
		panic(err)
	}

	writeAndRun(constant.EtcdSh, totalString)
}
func setFlannelConfig(targetIP string) {
	var (
		totalString string
		cmdLine     string
	)

	// 检查etcd运行状态
	cmdLine = "etcdctl  --endpoints http://" + targetIP + ":2379 member list\n"
	totalString += cmdLine + " \n"
	// 设置flannel参数
	cmdLine = "etcdctl set  /coreos.com/network/config '{\"Network\": \"10.0.0.0/16\", \"SubnetLen\": 24, \"SubnetMin\": \"10.0.10.0\",\"SubnetMax\": \"10.0.20.0\", \"Backend\": {\"Type\": \"vxlan\"}}'"
	totalString += cmdLine + " \n"

	writeAndRun(constant.TmpSh, totalString)
}

func startupFlannel(targetIP string) {
	var (
		totalString string
		cmdLine     string
	)

	// 启动flannel
	cmdLine = "flannel -etcd-endpoints \"http://" + targetIP + ":4001,http://" + targetIP + ":2379\""
	totalString += cmdLine + " &\n"

	writeAndRun(constant.FlannelSh, totalString)
}

func writeAndRun(filename string, context string) {
	var (
		get string
		err error
	)
	// 写入脚本
	script.Echo(context).WriteFile(filename)

	// 检验写入结果
	get, err = script.File(filename).String()

	if err != nil {
		panic(err)
	}

	fmt.Printf("check the file :\n %v", get)
	// 以下是执行该脚本内容的部分
	check("bash")

	goExecutable, _ := exec.LookPath("bash")

	cmdGoVer := &exec.Cmd{
		Path:   goExecutable,
		Args:   []string{goExecutable, filename},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	fmt.Println(cmdGoVer.String())

	if err := cmdGoVer.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
	os.Remove(constant.TmpSh)

}

func iptablesDelete(srcIP string, srcPort string, desIP string, desPort string) {
	var (
		totalString string
		cmdLine     string
	)
	cmdLine = "iptables -t nat -D OUTPUT -d " + srcIP + "/32 -p tcp -m tcp --dport " + srcPort + " -j DNAT --to-destination " + desIP + ":" + desPort

	totalString += cmdLine + " \n"
	cmdLine = "iptables -t nat -D PREROUTING -d " + srcIP + "/32 -p tcp -m tcp --dport " + srcPort + " -j DNAT --to-destination " + desIP + ":" + desPort
	totalString += cmdLine + " \n"
	writeAndRun(constant.TmpSh, totalString)
}

func iptablesSet(srcIP string, srcPort string, desIP string, desPort string) {
	var (
		totalString string
		cmdLine     string
	)
	cmdLine = "iptables -t nat -A OUTPUT -d " + srcIP + "/32 -p tcp -m tcp --dport " + srcPort + " -j DNAT --to-destination " + desIP + ":" + desPort

	totalString += cmdLine + " \n"
	cmdLine = "iptables -t nat -A PREROUTING -d " + srcIP + "/32 -p tcp -m tcp --dport " + srcPort + " -j DNAT --to-destination " + desIP + ":" + desPort
	totalString += cmdLine + " \n"
	writeAndRun(constant.TmpSh, totalString)
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

func GetPodIPByName(podName string) string {
	var (
		totalString string
		cmdLine     string
		err         error

		filename string = constant.TmpSh
		context  string
		res      string
	)
	cmdLine = "docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' " + podName + "-pause"
	totalString += cmdLine + " \n"

	context = totalString
	// 写入脚本
	script.Echo(context).WriteFile(filename)

	// 检验写入结果
	_, err = script.File(filename).String()

	if err != nil {
		panic(err)
	}

	//fmt.Printf("check the file :\n %v", get)
	// 以下是执行该脚本内容的部分
	//check("bash")
	var stdout, stderr bytes.Buffer

	goExecutable, _ := exec.LookPath("bash")

	cmdGoVer := &exec.Cmd{
		Path:   goExecutable,
		Args:   []string{goExecutable, filename},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	fmt.Println(cmdGoVer.String())

	if err := cmdGoVer.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
	res = stdout.String()
	//fmt.Printf("get stdout -> \n%v\n", stdout.String())
	//fmt.Printf("get stderr -> \n%v\n", stderr.String())
	os.Remove(constant.TmpSh)
	return res
}
