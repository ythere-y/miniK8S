package cmd

import (
	"fmt"
	"minik8s/apiserver"

	"github.com/spf13/cobra"
)

var RsCmd = &cobra.Command{
	Use:   "replicaset",
	Short: `replicaset 相关的命令`,
	Long:  `包括 get , update , delete , create`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(" 需要进一步指令,请查看 minik replicaset -h\n")
	},
}

var rscreate = &cobra.Command{
	Use:   "create",
	Short: "创建replicaset",
	Long:  "根据参数创建replicaset, 创建对应数量的Pod",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("create replicaset by file %v\n", args[0])
		apiserver.CreateRs(args[0])
	},
}

var rsdelete = &cobra.Command{
	Use:   "delete",
	Short: "删除replicaset",
	Long:  "根据指定的replicaset名字, 删除对应的replicaset",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("delete replicaset by name %v\n", args[0])
		apiserver.DeleteRs(args)
	},
}

var rsget = &cobra.Command{
	Use:   "get",
	Short: "获取replicaset的信息",
	Long:  "打印replicaset内的pod的各种信息",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("获取并打印replicaset内的pod信息")
		if len(args) == 0 {
			apiserver.DisplayAllRsInfo()
		} else {
			apiserver.DisplayRsInfo(args[0])
		}
	},
}

func init() {
	RsCmd.AddCommand(rscreate)
	RsCmd.AddCommand(rsdelete)
	RsCmd.AddCommand(rsget)
}
