package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/apiserver"
)

var TestCmd = &cobra.Command{
	Use:   "test",
	Short: `test 相关的命令`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("create pod by file %v\n", "podtest.yaml")
		apiserver.CmdCreatePod("podtest.yaml")
	},
}
