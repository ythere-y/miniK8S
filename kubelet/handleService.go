package kubelet

import (
	"encoding/json"
	"fmt"
	"github.com/bitfield/script"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"minik8s/apiserver"
	"minik8s/constant"
	. "minik8s/lab/etcd"
	"minik8s/registry/service"
	"os"
	"os/exec"
)

func kubeleteServiceWatch() {

	watchName := SetKey(
		SetPrefix(constant.RegistryPrefix),
		SetSourceType(constant.ServiceSourceName),
	)
	apiserver.SyncWatch(watchName, servicehandler)

}
func servicehandler(event *clientv3.Event) error {
	var err error

	//serviceName := utils.GetLastWord(string(event.Kv.Key))

	switch event.Type {

	case mvccpb.DELETE:
		// 删除一个service的操作
		fmt.Printf("kubelet handling Delete service = %v\n", string(event.Kv.Key))

		//TODO 有点问题，暂时不写

	case mvccpb.PUT:
		// 增加一个services的操作
		fmt.Printf("kubelet handling key = %v, value = %v\n", string(event.Kv.Key), string(event.Kv.Value))
		var servieInfo service.Service
		err = json.Unmarshal(event.Kv.Value, &servieInfo)

	}
	return err
}

func AppendIptables(ser service.Service) {

}

func iptables_A(srcIP string, Port string, desIP string) {
	var (
		cmdLine     string
		totalString string
		get         string
		err         error
	)

	cmdLine = ""
	totalString = ""

	cmdLine = "iptables -t nat -D OUTPUT -d " + srcIP + ":" + Port + " -p tcp -m tcp --dport 10 -j DNAT --to-destination " + desIP + ":" + Port
	totalString += cmdLine

	cmdLine = "iptables -t nat -D PREROUTING -d " + srcIP + ":" + Port + " -p tcp -m tcp --dport 10 -j DNAT --to-destination " + desIP + ":" + Port
	totalString += cmdLine

	// 写入脚本
	script.Echo(totalString).WriteFile(constant.EtcdSh)

	// 检验写入结果
	get, err = script.File(constant.EtcdSh).String()
	if err != nil {
		panic(err)
	}
	fmt.Printf("check the file :\n %v", get)

	// 以下是执行该脚本内容的部分

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

func iptables_D(srcIP string, Port string, desIP string) {
	var (
		cmdLine     string
		totalString string
		get         string
		err         error
	)

	cmdLine = ""
	totalString = ""

	cmdLine = "iptables -t nat -D OUTPUT -d " + srcIP + ":" + Port + " -p tcp -m tcp --dport 10 -j DNAT --to-destination " + desIP + ":" + Port
	totalString += cmdLine

	cmdLine = "iptables -t nat -D PREROUTING -d " + srcIP + ":" + Port + " -p tcp -m tcp --dport 10 -j DNAT --to-destination " + desIP + ":" + Port
	totalString += cmdLine

	// 写入脚本
	script.Echo(totalString).WriteFile(constant.EtcdSh)

	// 检验写入结果
	get, err = script.File(constant.EtcdSh).String()
	if err != nil {
		panic(err)
	}
	fmt.Printf("check the file :\n %v", get)

	// 以下是执行该脚本内容的部分

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
