package main

import (
	"minik8s/service"
)

var rootName string

func main() {
	//cmd.RootCmdRun()// 关于命令行的测试

	service.SerReadTest() // 关于读取yamle文件建立service的测试

	//circle.CircleTest() // 关于循环import的测试
	return

}
