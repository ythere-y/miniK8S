package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "关闭",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(1)
	},
}
