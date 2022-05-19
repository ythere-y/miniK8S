package utils

import (
	"fmt"
	"runtime"
	"strings"
)

func FirstWord(str string) string {
	//存放单词的切片
	var strSlice []string
	//以空格分隔句子
	strSlice = strings.Split(str, " ")
	return strSlice[0]
}

func CutFirst(str string, rootname string) string {
	str = strings.Replace(str, rootname, " ", 1)
	return str
}
func HandleError(words string, err error) {
	if err != nil {
		fmt.Printf("%v ->: \n%v\n", words, err.Error())
	}
}

func CheckNil(name string, v any) {
	if v == nil {
		DebugInfo(name + " is nil !")
	} else {
		DebugInfo(name + " is NOT nil")
	}
}

func DebugInfo(str string) {
	fmt.Println("[DEBUG] [INFO] " + str)

}
func DebugError(str string) {
	fmt.Println("[DEBUG] [ERROR] " + str)
}
func DebugTanInfo() {
	fmt.Println("get here")
	return
	funcName, file, line, ok := runtime.Caller(0)
	if ok {
		fmt.Println("func name: " + runtime.FuncForPC(funcName).Name())
		fmt.Printf("file: %s, line: %d\n", file, line)
	}
}
