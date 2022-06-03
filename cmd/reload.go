package cmd

import (
	"github.com/spf13/cobra"
)

var ReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: `容错相关的命令`,
	Long:  `集群重载命令`,
	Run: func(cmd *cobra.Command, args []string) {
		//TODO:容错，目前思路是直接重启master和worker，因为其他信息都持久化在etcd

	},
}
