package api

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func TestFuncDownload(t *testing.T) {
	sshConfig := &ssh.ClientConfig{
		User: "ossftp",
		Auth: []ssh.AuthMethod{
			ssh.Password("oss123"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// SSH 连接
	sshClient, err := ssh.Dial("tcp", "192.168.33.10:2222", sshConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer sshClient.Close()

	// 创建 SFTP 客户端
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		log.Fatal(err)
	}
	defer sftpClient.Close()
	// 这里的/才是关键所\var\opt\ericsson\GLPGW1\aaa.xml
	// 是windows系统跑的代码和linux系统跑的代码不一致导致的
	remoteFilePath := "/var/opt/ericsson/GLPGW1/aaa.xml"
	remoteFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer remoteFile.Close()
}

func TestCreateFolder(t *testing.T) {
	tmpPath := "C:\\workspace\\Go_download\\download"
	folderName := "GLCUDB1"
	targetDir := filepath.Join(tmpPath, folderName)
	_, err := os.Stat(targetDir)
	if os.IsNotExist(err) {
		// 文件夹不存在，创建文件夹
		err := os.Mkdir(targetDir, 0755)
		if err != nil {
			fmt.Println("无法创建文件夹:", err)
		}
		fmt.Println("文件夹已创建")
	}
}

func TestTime(t *testing.T) {
	myStr := `2023-05-31T10:00`
	layout := "2006-01-02T15:04" // 定义日期时间格式
	ta, err := time.Parse(layout, myStr)
	if err != nil {
		fmt.Println(err)
	}
	currnetTime := time.Now()
	value := currnetTime.Sub(ta)
	fmt.Println("curent ", currnetTime)
	fmt.Println("ta", ta)
	fmt.Println(value)
}
