package cmd

import (
	"bufio"
	"fmt"
	"minik8s/utils"
	"os"
	"strings"

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
	RootCmd.AddCommand(TestCmd)
	RootCmd.AddCommand(joinCmd)
	RootCmd.AddCommand(RsCmd)
	RootCmd.AddCommand(AutoScalerCmd)
	RootCmd.AddCommand(ReloadCmd)
	RootCmd.AddCommand(startCmd)
	RootCmd.AddCommand(NodeCmd)
	RootCmd.AddCommand(exitCmd)
}

var rootName string

func RootCmdRun() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please input your command, type 'quit' to exit.")
	for true {
		fmt.Print("-> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(fmt.Errorf("fatal error: %w \n", err))
		}
		input = strings.TrimRight(input, "\r\n")
		if input == "quit" {
			break
		}
		rootName = RootCmd.Name()
		if strings.Compare(utils.FirstWord(input), rootName) == 0 {
			input = utils.CutFirst(input, rootName)
			fmt.Println("Executing command ...")

			RootCmd.SetArgs(strings.Fields(input))

			err = RootCmd.Execute()
		} else {
			fmt.Printf("shold start with %v\n", rootName)
		}

	}
}
