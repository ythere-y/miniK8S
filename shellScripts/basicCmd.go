package shellScripts

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func check(name string) string {
	goExecPath, err := exec.LookPath(name)

	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Go Executable: ", goExecPath)
	}
	return goExecPath

}

func gocmdTest() {
	check("go")
	goExecutable, _ := exec.LookPath("go")

	cmdGoVer := &exec.Cmd{
		Path:   goExecutable,
		Args:   []string{goExecutable, "version"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	fmt.Println(cmdGoVer.String())

	if err := cmdGoVer.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
}
func sleepTest() {
	name := "sleep"
	check(name)
	goExecutable, _ := exec.LookPath(name)
	cmd := &exec.Cmd{
		Path:   goExecutable,
		Args:   []string{goExecutable, "3"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	for i := 1; i < 3000; i++ {
		fmt.Println(i)
	}
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}

func goVersionTest() {
	check("go")
	cmd := exec.Command("sleep", "60")
	var outb, errb bytes.Buffer

	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Printf("sleep end")
	//cmd.Wait()
}
