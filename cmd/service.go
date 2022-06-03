package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/apiserver"
)

var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: `service 相关的命令`,
	Long:  `包括 get , delete , create`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("service 需要进一步指令,请查看 minik service -h\n")
	},
}

var serviceget = &cobra.Command{
	Use:   "get",
	Short: "获取service的信息",
	Long:  "打印service的各种信息, get / get name",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			apiserver.DisplayAllServiceInfo()
		} else {
			apiserver.DisplayNodeInfo(args[0])
		}
	},
}

var servicedelete = &cobra.Command{
	Use:   "delete",
	Short: "删除service",
	Long:  "可以删除service, delete name1 name2 ...",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiserver.CmdDeleteService(args)
	},
}

var servicecreate = &cobra.Command{
	Use:   "create",
	Short: "获取service的信息",
	Long:  "打印service的各种信息, create file.name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiserver.CreateServiceByFile(args[0])
	},
}

func init() {
	ServiceCmd.AddCommand(servicecreate)
	ServiceCmd.AddCommand(serviceget)
	ServiceCmd.AddCommand(servicedelete)
}
