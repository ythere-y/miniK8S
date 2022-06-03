package utils

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	yaml "gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

var Dir = "D:/Schoolwork/2022_spring/CloudComputing/labs/Minik8s/minik8s/cuda/"

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

	//var stdOut, stdErr bytes.Buffer

	//session, err := SSHConnect("stu610", "C3Dr6MM%", "login.hpc.sjtu.edu.cn", 22)
	session, err := SSHConnect("stu610", "C3Dr6MM%", "login.hpc.sjtu.edu.cn", 22)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer session.Close()

	//session.Stdout = &stdOut
	//session.Stderr = &stdErr
	//
	//session.Run("if [ -e /lustre ]; then echo 0; else echo 1; fi")
	//session.Run("cd")
	//session.Run("mkdir aaa")
	//ret, err := strconv.Atoi(str.Replace(stdOut.String(), "\n", "", -1))
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("%d, %s\n", ret, stdErr.String())

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
	sftpClient, err = sftpconnect("stu610", "C3Dr6MM%", "login.hpc.sjtu.edu.cn", 22)
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

func CopyandRun() error {
	err := ScpCopy(Dir+"test.slurm", "/lustre/home/acct-stu/stu610/cuda")
	check(err)
	err = ScpCopy(Dir+"test002.cu", "/lustre/home/acct-stu/stu610/cuda")
	check(err)
	RunSsh("whoami")
	return err
}

func ParseYaml(file string) Slurm {
	fmt.Println("start parsing yaml file")
	var slurm Slurm
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("yaml file read error")
		fmt.Println(err)
	}

	err = yaml.Unmarshal(yamlFile, &slurm)
	if err != nil {
		fmt.Println("yaml unmarshal error")
		fmt.Println(err)
	}
	return slurm
}

func Slu2File(slurm Slurm) string {
	fmt.Println("start generating slurm file")
	var filename = Dir + slurm.MetaData.Name + ".slurm"
	var writeString = "#!/bin/bash\n\n" +
		"#SBATCH --job-name=" + slurm.MetaData.Name +
		"\n#SBATCH --partition=" + slurm.MetaData.Partition +
		"\n#SBATCH --output=" + slurm.Spec.Output +
		"\n#SBATCH --error=" + slurm.Spec.Error +
		"\n#SBATCH -N " + slurm.Spec.N +
		"\n#SBATCH --ntasks-per-node=" + slurm.Spec.TaskPerNode +
		"\n#SBATCH --cpus-per-task=" + slurm.Spec.CpuPerTask +
		"\n#SBATCH --gres=gpu:" + slurm.Spec.Gres.Gpu + "\n\n" +
		"ulimit -s unlimited\n" +
		"ulimit -l unlimited\n\n" +
		"module load gcc/8.3.0 cuda/10.1.243-gcc-8.3.0\n" +
		"nvcc " + slurm.MetaData.Name + ".cu -o " + slurm.MetaData.Name + " -lcublas\n\n" +
		"./" + slurm.MetaData.Name

	var f *os.File
	var err1 error

	if checkFileIsExist(filename) { //如果文件存在
		err1 = os.Remove(filename)
		f, err1 = os.Create(filename) //创建文件
		//f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
		fmt.Println("文件存在")
	} else {
		f, err1 = os.Create(filename) //创建文件
		fmt.Println("文件不存在")
	}
	check(err1)
	n, err1 := io.WriteString(f, writeString) //写入文件(字符串)
	check(err1)
	fmt.Printf("写入 %d 个字节\n", n)
	return filename
}

func Submit(yfile string, cfile string) {
	slurmFile := Slu2File(ParseYaml(yfile))
	err := ScpCopy(slurmFile, "/lustre/home/acct-stu/stu610/cuda")
	check(err)
	err = ScpCopy(cfile, "/lustre/home/acct-stu/stu610/cuda")
	check(err)
	ret := RunSsh("cd /lustre/home/acct-stu/stu610/cuda; sbatch " + path.Base(slurmFile))
	fmt.Println(ret)

}

func GetStat(name string) []string {
	res := RunSsh("sacct | grep " + name)
	arr := strings.Fields(res)
	size := len(arr)
	fmt.Println("JobId: " + arr[size-7] + "\t\tName: " + arr[size-6] + "\t\tState: " + arr[size-2])
	return arr
}

func GetRes(name string) {
	res := RunSsh("sacct | grep " + name)
	arr := strings.Fields(res)
	size := len(arr)
	if arr[size-2] == "COMPLETED" || arr[size-2] == "FAILED" {
		str := RunSsh("cd cuda; cat " + arr[size-7] + ".out")
		fmt.Println("output:\n" + str)
	} else {
		fmt.Println("Job uncompleted\n" + "State: " + arr[size-2])
	}

}

type Slurm struct {
	Kind     string `yaml:"kind"`
	MetaData struct {
		Name      string `yaml:"name"`
		Partition string `yaml:"partition"`
	}
	Spec struct {
		N           string `yaml:"n"`
		TaskPerNode string `yaml:"task-per-node"`
		CpuPerTask  string `yaml:"cpu-per-task"`
		Gres        struct {
			Gpu string `yaml:"gpu"`
		}
		Output string `yaml:"output"`
		Error  string `yaml:"error"`
	}
}
