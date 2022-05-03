package lab

import (
	"fmt"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
func init() {

	//registry.RootCmd.AddCommand(EchoCmd)
	fmt.Println("hello init")
}

func Cmd() *cobra.Command {
	nodeCmd.AddCommand(startCmd())
	return nodeCmd
}

var EchoCmd = &cobra.Command{
	// 子命令描述，第一个单词作为子命令名称
	Use: "echo",

	// 子命令别名，可以使用任何一个别名替代上面的 echo 来运行子命令
	Aliases: []string{
		"copy",
		"repeat",
	},

	// 命令提示，默认情况如果输入的子命令不存在会提示 unknown cmd，
	// 但是如果定义了 SuggestFor 的情况下，如果输入的命令不存在，会去 SuggestFor
	// 里面查找是否有匹配的字符串，如果有，则提示是否期望输入的是 echo 命令
	SuggestFor: []string{
		"ech", "cp", "rp",
	},

	DisableSuggestions:         true,
	SuggestionsMinimumDistance: 2,
	DisableFlagsInUseLine:      true,

	// 简单扼要概括下命令的用途
	Short: "echo is a command to echo command line inputs",

	// 想说什么都在这里说吧，越详细越好，可以使用 `` 来跨行输入
	Long: `echo is a command to echo command line inputs.
It is a very simple command used to display how to implement command line tools
using cobra, which is a very famous library to build command line interface tools.
`,

	// 新版本废弃时隐藏，但是实际是可以执行的
	//Hidden: true,

	// 废弃说明
	//Deprecated: "will be deleted in version 2.0",

	// 注解，用于代码层面的命令分组，不会显示在命令行输出中
	Annotations: map[string]string{
		"group":        "user",
		"require-auth": "none",
	},

	SilenceErrors: true,
	//SilenceUsage:  true,
	// Version 定义版本
	Version: "1.0.0",

	// 是否禁用选项解析
	//DisableFlagParsing: true,

	//DisableAutoGenTag: true,

	// PersistentPreRun: func(cmd *cobra.Command, args []string) {
	//     fmt.Println("let me check echk")
	// },

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		//    return errors.New("invalid parameter")
		return nil
	},

	// 子命令执行过程
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("RunE")
		return nil
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("PostRun")
	},
}

const nodeFuncName = "node"

var (
	stopPidFile string
)
var nodeCmd = &cobra.Command{
	Use:   nodeFuncName,
	Short: fmt.Sprintf("%s specific commands.", nodeFuncName),
	Long:  fmt.Sprintf("%s specific commands.", nodeFuncName),
}

var chaincodeDevMode bool

func startCmd() *cobra.Command {
	return nodeStartCmd
}

var nodeStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the node.",
	Long:  `Starts a node that interacts with the network.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return serve(args)
	},
}

func serve(args []string) error {
	fmt.Println("不要企图爱上哥，哥只是个传说:", args)
	return nil
}
