package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "demo",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//    Run: func(cmd *cobra.Command, args []string) { },
}

func init() {
	RootCmd.AddCommand(addCmd)

	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(serveCmd)
}
