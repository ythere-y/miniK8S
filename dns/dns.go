package dns

import (
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func Ip2name(ip string, name string) {
	if checkFileIsExist("/etc/hosts") {
		file, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		if _, err = file.WriteString(ip + " " + name + "\n"); err != nil {
			panic(err)
		}
	} else {
		fmt.Println("file /etc/hosts does not exist")
		panic(0)
	}

	fmt.Println("Domain name config success")
}
