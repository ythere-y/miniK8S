package utils

import (
	"fmt"
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
