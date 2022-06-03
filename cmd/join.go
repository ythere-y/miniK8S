package cmd

import (
	"github.com/spf13/cobra"
	"minik8s/kubernetes"
)

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "加入一个master，成为一个worker",
	Long:  `join masterIP:port node.yaml`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		kubernetes.JoinAsWorker(args)
	},
}
