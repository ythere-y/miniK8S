package dns

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

func ParseYaml(file string) Dns {
	fmt.Println("start parsing yaml file")
	var dns Dns
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("yaml file read error")
		fmt.Println(err)
	}

	err = yaml.Unmarshal(yamlFile, &dns)
	if err != nil {
		fmt.Println("yaml unmarshal error")
		fmt.Println(err)
	}
	return dns
}

func C2host(dns Dns) {
	for _, val := range dns.DnsBinds {
		Ip2name(val.ServiceIp, val.Path+"."+dns.Host)
	}
}

type Dns struct {
	Kind     string    `yaml:"kind"`
	Name     string    `yaml:"name"`
	Host     string    `yaml:"host"`
	DnsBinds []DnsBind `yaml:"binds"`
}

type DnsBind struct {
	ServiceIp string `yaml:"service-ip"`
	Path      string `yaml:"path"`
}
