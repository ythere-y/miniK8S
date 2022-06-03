package cmd

import (
	"github.com/spf13/cobra"
	"minik8s/kubernetes"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "开启集群，成为一个master",
	Long:  `start, 需要配置文件 ./config/masterNode.yaml`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		kubernetes.StartUpMaster()
	},
}
