package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/service"
)

var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: `service 相关的命令`,
	Long:  `包括 get , update , delete , create`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("service 需要进一步指令,请查看 minik service -h\n")
	},
}

var serviceget = &cobra.Command{
	Use:   "get",
	Short: "获取service的信息",
	Long:  "打印service的各种信息",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service.ServicePreDisplay()
		service.GetServiceByName(args[0])
	},
}

var servicedelete = &cobra.Command{
	Use:   "delete",
	Short: "删除service",
	Long:  "可以删除service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service.DeleteServiceByName(args[0])
	},
}

var servicecreate = &cobra.Command{
	Use:   "create",
	Short: "获取service的信息",
	Long:  "打印service的各种信息",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service.CreateServcieByFile(args[0])
	},
}

var serviceupdate = &cobra.Command{
	Use:   "update",
	Short: "更新service的信息",
	Long:  "更新service的各种信息",
	Run: func(cmd *cobra.Command, args []string) {
		//TODO:需要接上正确的接口
		fmt.Println("update fixing")
	},
}

func init() {
	ServiceCmd.AddCommand(servicecreate)
	ServiceCmd.AddCommand(serviceget)
	ServiceCmd.AddCommand(serviceupdate)
	ServiceCmd.AddCommand(servicedelete)
}
