package main

import (
	"bufio"
	"fmt"
	"github.com/docker/docker/client"
	"minik8s/lab/cmd"
	"minik8s/lab/dksdk"
	"minik8s/utils"
	"os"
	"strings"
	"time"
)

var rootName string

func main() {
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
		rootName = cmd.RootCmd.Name()
		if strings.Compare(utils.FirstWord(input), rootName) == 0 {
			input = utils.CutFirst(input, rootName)
			fmt.Println("Executing command ...")

			cmd.RootCmd.SetArgs(strings.Fields(input))

			err = cmd.RootCmd.Execute()
		} else {
			fmt.Printf("shold start with %v\n", rootName)
		}

	}

	return

	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
	}
	id := dksdk.CreateContainer(cli, "library/alpine", []string{"echo", "hello world"}, "testCreate", nil, nil)
	fmt.Printf("%s\n", id)
	time.Sleep(time.Second * 1)
	dksdk.StartContainer(id, cli)
	time.Sleep(time.Second * 40)
	dksdk.StopContainer(id, cli)
	time.Sleep(time.Second * 3)
	id, err = dksdk.RemoveContainer(id, cli)
	if err == nil {
		fmt.Println("删除容器", id, "成功")
	}
}
