package cmd

import (
	"github.com/spf13/cobra"
	"minik8s/dns"
)

var DnsCmd = &cobra.Command{
	Use:   "dns",
	Short: `DNS 配置命令`,
	Long:  `DNS 配置命令`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dns.Run(args[0])
	},
}
