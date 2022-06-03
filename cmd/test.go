package cmd

import (
	"github.com/spf13/cobra"
	"minik8s/apiserver"
	"minik8s/kubernetes"
)

var TestCmd = &cobra.Command{
	Use:   "test",
	Short: `test 相关的命令`,
	Run: func(cmd *cobra.Command, args []string) {

		kubernetes.StartUpMaster()
		//apiserver.CmdCreatePod("podtest.yaml")
		apiserver.CreateServiceByFile("servicetest.yaml")
		//apiserver.CmdDeleteNode([]string{"testpod"})
	},
}
