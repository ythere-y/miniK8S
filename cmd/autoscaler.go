package cmd

import (
	"fmt"
	"minik8s/apiserver"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var AutoScalerCmd = &cobra.Command{
	Use:   "autoscaler",
	Short: `autoscaler 相关的命令`,
	Long:  `包括 get , update , delete , create`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("autoscaler 需要进一步指令,请查看 minik autoscaler -h\n")
	},
}

var asget = &cobra.Command{
	Use:   "get",
	Short: "获取autoscaler的信息",
	Long:  "打印autoscaler的各种信息",
	Args:  cobra.MaximumNArgs(1),
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("获取并打印autoscaler信息")
	// 	if len(args) == 0 {
	// 		apiserver.DisplayAllPodsInfo()
	// 	} else {
	// 		apiserver.DisplayPodsInfo(args[0])
	// 	}
	// },
}

var asdelete = &cobra.Command{
	Use:   "delete",
	Short: "删除autoscaler",
	Long:  "可以删除autoscaler",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("delete autoscaler by name %v\n", args[0])
		apiserver.DeletePod(args)
	},
}

// var asstop = &cobra.Command{
// 	Use:   "stop",
// 	Short: "停止pod",
// 	Long:  "可以停止pod的运行，但是pod仍然属于对应的node。只是将pod管理下的所有container状态设置为stopped",
// 	Args:  cobra.MinimumNArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Printf("stop pod by name %v\n", args[0])
// 		apiserver.StopPod(args)
// 	},
// }

// minik autoscaler create astest.yaml
var ascreate = &cobra.Command{
	Use:   "create",
	Short: "创建autoscaler",
	Long:  "根据参数创建autoscaler",
	Args:  cobra.ExactArgs(1),
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Printf("create pod by file %v\n", args[0])
	// 	apiserver.CreatePod(args[0])
	// },
}

func args(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {

		return errors.New("至少需要三个参数!")
	}
	return nil
}

var asupdate = &cobra.Command{
	Use:   "update",
	Short: "更新autoscaler的信息",
	Long:  "更新autoscaler的各种信息",
	// Run: func(cmd *cobra.Command, args []string) {
	// 	//TODO:需要接上正确的接口
	// 	fmt.Println("update fixing")
	// },
}

func init() {
	PodCmd.AddCommand(ascreate)
	PodCmd.AddCommand(asget)
	PodCmd.AddCommand(asupdate)
	PodCmd.AddCommand(asdelete)
}
