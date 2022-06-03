package dns

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

//minik dns ./dns/dns.yaml

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

func C2Rhost(dns Dns) {
	//ScpCopy("./dns/dns.yaml", "/tmp/minik/minik8s/dns")
	for _, val := range dns.DnsBinds {
		RunSsh("cat >> /etc/hosts << EOF\n" + val.ServiceIp + " " + val.Path + "." + dns.Host + "\nEOF")
	}
}

func RunCmd(cmd string) string {
	runCmd := cmd
	_, err := script.Echo(runCmd).WriteFile("./dns/run.sh")
	if err != nil {
		panic(err)
	}
	_, err = script.File("./dns/run.sh").String()
	//fmt.Printf("check the file :\n %v", get)
	Path := check("bash")
	cmdGoVer := &exec.Cmd{
		Path: Path,
		Args: []string{Path, "./dns/run.sh"},
		//Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	///fmt.Println("OUT", cmdGoVer.String())
	out, err := cmdGoVer.Output()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	//fmt.Println(string(out))
	//str := strings.Replace(string(out), "\n", "", -1)
	str := string(out)
	return str
}

func C2Dock(dns Dns) {
	res := RunCmd("docker ps")
	arr := strings.Split(res, "\n")
	var names []string
	for i, val := range arr {
		if i == 0 {
			continue
		}

		line := strings.Fields(val)
		if len(line) == 0 {
			continue
		}
		names = append(names, line[len(line)-1])
	}
	fmt.Println(names)
	for _, name := range names {
		for _, bind := range dns.DnsBinds {
			cmd := "docker exec -i " + name + " /bin/sh -c \"cat >> /etc/hosts << EOF\n" + bind.ServiceIp + " " + bind.Path + "." + dns.Host + "\nEOF\""
			RunCmd(cmd)
		}
	}

}

func check(name string) string {
	goExecPath, err := exec.LookPath(name)

	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Go Executable: ", goExecPath)
	}

	return goExecPath

}

func Run(file string) {
	out := ParseYaml("./dns/dns.yaml")
	C2host(out)
	C2Rhost(out)
	C2Dock(out)
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

func SSHConnect(user, password, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		// Timeout:             30 * time.Second,
		HostKeyCallback: hostKeyCallbk,
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}

func RunSsh(cmd string) string {

	session, err := SSHConnect("root", "helloHC123", "10.119.11.73", 22)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer session.Close()

	combo, err := session.CombinedOutput(cmd)

	if err != nil {
		log.Fatal("远程执行cmd 失败", err)

	}

	//log.Println("命令输出:", string(combo))
	return string(combo)

}

func sftpconnect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		//这个是问你要不要验证远程主机，以保证安全性。这里不验证
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

//单个copy
func ScpCopy(localFilePath, remoteDir string) error {
	var (
		sftpClient *sftp.Client
		err        error
	)
	// 这里换成实际的 SSH 连接的 用户名，密码，主机名或IP，SSH端口
	sftpClient, err = sftpconnect("root", "helloHC123", "10.119.11.73", 22)
	if err != nil {
		log.Println("scpCopy:", err)
		return err
	}
	defer sftpClient.Close()
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		log.Println("scpCopy:", err)
		return err
	}
	defer srcFile.Close()

	var remoteFileName = path.Base(localFilePath)
	dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		log.Println("scpCopy:", err)
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf[0:n])
	}
	return nil
}
