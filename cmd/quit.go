package cmd

import (
	"encoding/json"
	"minik8s/apiserver"
	"minik8s/constant"
	"minik8s/lab/etcd"
	"minik8s/registry/pod"
	"minik8s/utils"
	"os"
	"time"

	"github.com/spf13/cobra"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "关闭",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			num    uint32 = 0
			getRsp *clientv3.GetResponse
			err    error
		)
		getRsp, err = etcd.GetWithPrefix(
			etcd.SetKey(
				etcd.SetPrefix(constant.RegistryPrefix),
				etcd.SetSourceType(constant.PodSourceName)))
		utils.HandleError("get with prefix error[from get child num]", err)
		num = uint32(len(getRsp.Kvs))
		if num == 0 {
			return
		}
		var names []string // get all pods name
		for _, event := range getRsp.Kvs {
			tmpValue := event.Value
			var tmpPod pod.Pod
			err = json.Unmarshal(tmpValue, &tmpPod)
			podName := tmpPod.Meta.Name
			utils.HandleError("unmarshal pod error", err)
			names = append(names, podName)
		}
		apiserver.CmdStopPods(names)
		time.Sleep(1 * time.Second)
		os.Exit(1)
	},
}
