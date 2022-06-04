package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/utils"
)

var CudaCmd = &cobra.Command{
	Use:   "cu",
	Short: `GPU 相关的命令`,
	Long:  `包括 submit , state , result`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("需要进一步指令,请查看 minik cu -h\n")
	},
}

//minik cu submit D:/Schoolwork/2022_spring/CloudComputing/labs/Minik8s/minik8s/cuda/test.yaml D:/Schoolwork/2022_spring/CloudComputing/labs/Minik8s/minik8s/cuda/matrix.cu
var CudaSubmit = &cobra.Command{
	Use:   "submit",
	Short: "提交任务",
	Long:  "提交任务到交我算平台并编译运行",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		utils.Submit(args[0], args[1])
	},
}

//minik cu state matrix
var CudaState = &cobra.Command{
	Use:   "state",
	Short: "获取任务状态",
	Long:  "获取任务状态 usage: minik cu state JobName",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utils.GetStat(args[0])
	},
}

//minik cu result matrix
var CudaRes = &cobra.Command{
	Use:   "result",
	Short: "获取任务结果",
	Long:  "获取任务结果 usage: minik cu result JobName",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utils.GetRes(args[0])
	},
}

func init() {
	CudaCmd.AddCommand(CudaSubmit)
	CudaCmd.AddCommand(CudaState)
	CudaCmd.AddCommand(CudaRes)

}
