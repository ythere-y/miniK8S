package shellScripts

import (
	"fmt"
	"github.com/bitfield/script"
	"minik8s/constant"
	"os"
	"os/exec"
)

//readSubnetEnv
/*
自动读取/run/flannel/subnet.env的内容

取出其中需要用于docker参数的内容

第一个string是IP

第二个string是MTU

第三个返回值是一个恒为nil的error
*/
func readSubnetEnv() (string, string, error) {
	var (
		getIP   string
		getMTU  string
		cutFlag bool = false
		get     string
		err     error
	)
	get, err = script.File("/run/flannel/subnet.env").Match("FLANNEL_SUBNET").String()
	if err != nil {
		panic(err)
	}
	cutFlag = false
	for i := 0; i < len(get); i++ {
		if cutFlag {
			getIP += string(get[i])
			continue
		}
		if get[i] >= '0' && get[i] <= '9' {
			cutFlag = true
			i--
		}
	}
	get, err = script.File("/run/flannel/subnet.env").Match("FLANNEL_MTU").String()

	cutFlag = false
	for i := 0; i < len(get); i++ {
		if cutFlag {
			getMTU += string(get[i])
			continue
		}
		if get[i] >= '0' && get[i] <= '9' {
			cutFlag = true
			i--
		}
	}

	fmt.Printf("read the /run/flannel/subnet.env :\n%v", get)
	fmt.Printf("get IP = %v\n", getIP)
	fmt.Printf("get MTU = %v\n", getMTU)

	return getIP, getMTU, nil
}

func BuildDockerCommandFix() (string, error) {
	var (
		ip  string
		mtu string
		err error
		ret string = ""
	)

	ip, mtu, err = readSubnetEnv()
	if err != nil {
		panic(err)
	}
	ret += " --bip=" + ip
	ret += " --ip-masq=true"
	ret += " --mtu=" + mtu

	fmt.Printf("the fix = %v\n", ret)
	return ret, nil
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

// EtcdSetUp
/*
环境设置，启动etcd
*/
func EtcdSetUp(thisip string) {
	startupEtcd(thisip)
}

// FlannelSetUp
/*
启动flannel,需要目标的ip
*/
func FlannelSetUp(targetip string) {
	startupFlannel(targetip)
}

func MasterEnvSetUp(ip string) {
	EtcdSetUp(ip)
	FlannelSetUp(ip)
}

func SlaveEnvSetUp(thisip string, masterip string) {
	FlannelSetUp(masterip)
}
