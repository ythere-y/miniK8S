package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "minik",
	Short: "minik8s命令行系统",
	Long:  `可以通过命令行和minik8s系统交互`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//    Run: func(cmd *cobra.Command, args []string) { },
}

func init() {
	RootCmd.AddCommand(addCmd)
	RootCmd.AddCommand(ServiceCmd)
	RootCmd.AddCommand(PodCmd)
}
