package cmd

import "github.com/spf13/cobra"

var serveCmd = &cobra.Command{
	Use: "serve",
	Short: `短的描述, 短的描述在不指定使用这个命令时调用：
				例如：newApp -h`,
	Long: `长的描述`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
