package environment

import (
	"fmt"
	"os/exec"
	"time"
)

func Main() {
	var (
		command *exec.Cmd
		err     error
	)
	//command = exec.Command("cd", "'E:\\'")
	//err = command.Run()
	//if err != nil {
	//	panic(err)
	//}

	command = exec.Command("notepad")
	err = command.Run()
	if err != nil {
		panic(err)
	}

	command = exec.Command("pwd")
	err = command.Run()
	if err != nil {
		panic(err)
	}

	tick := time.Tick(time.Second)
	for range tick {
		fmt.Println("...")
	}
}
