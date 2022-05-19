package cmd

import (
	"minik8s/apiserver"
	"minik8s/replicaset"

	"github.com/spf13/cobra"
)

var RsCmd &cobra.Command {
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
		fmt.Printf("create pod by file %v\n", args[0])
		apiserver.CreateRs(args[0])
	},
}

func args(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("至少需要三个参数!")
	}
	return nil
}

func init() {
	RsCmd.AddCommand(rscreate)
}