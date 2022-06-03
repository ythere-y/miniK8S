package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/apiserver"
)

var NodeCmd = &cobra.Command{
	Use:   "node",
	Short: `node 相关的命令`,
	Long:  `包括 get , delete`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("node 需要进一步指令,请查看 minik node -h\n")
	},
}

var nodeget = &cobra.Command{
	Use:   "get",
	Short: "获取node的信息",
	Long:  "打印node的各种信息, node / node name",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("获取并打印node信息")
		if len(args) == 0 {
			apiserver.DisplayAllNodeInfo()
		} else {
			apiserver.DisplayNodeInfo(args[0])
		}
	},
}

var nodedelete = &cobra.Command{
	Use:   "delete",
	Short: "删除node",
	Long:  "可以删除node, delete name1 name2 ...",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("delete node by name %v\n", args[0])
		apiserver.CmdDeleteNode(args)
	},
}

func init() {
	NodeCmd.AddCommand(nodeget)
	NodeCmd.AddCommand(nodedelete)
}
