package cmd

import (
	"fmt"
	"minik8s/apiserver"
	"minik8s/pod"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var PodCmd = &cobra.Command{
	Use:   "pod",
	Short: `pod 相关的命令`,
	Long:  `包括 get , update , delete , create`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("service 需要进一步指令,请查看 minik pod -h\n")
	},
}

var podget = &cobra.Command{
	Use:   "get",
	Short: "获取pod的信息",
	Long:  "打印pod的各种信息",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("获取并打印pod信息")
		pod.PodPreDisplay()
		pod.GetPodInfoByName(args[0])
	},
}

var poddelete = &cobra.Command{
	Use:   "delete",
	Short: "删除pod",
	Long:  "可以删除pod",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("delete pod by name %v\n", args[0])
		pod.StopPodByName(args[0])
		pod.RemovePodByName(args[0])
	},
}

// minik pod create file.yaml
var podcreate = &cobra.Command{
	Use:   "create",
	Short: "创建pod",
	Long:  "根据参数创建pod",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("create pod by file %v\n", args[0])
		apiserver.CreatePod(args[0])
	},
}

func args(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {

		return errors.New("至少需要三个参数!")
	}
	return nil
}

var podupdate = &cobra.Command{
	Use:   "update",
	Short: "更新pod的信息",
	Long:  "更新pod的各种信息",
	Run: func(cmd *cobra.Command, args []string) {
		//TODO:需要接上正确的接口
		fmt.Println("update fixing")
	},
}

func init() {
	PodCmd.AddCommand(podcreate)
	PodCmd.AddCommand(podget)
	PodCmd.AddCommand(podupdate)
	PodCmd.AddCommand(poddelete)
}
