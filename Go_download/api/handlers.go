package api

import (
	"fmt"
	config "godownload/configs"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func HelloHandler(c *gin.Context) {
	startTimeStr := c.Query("start")
	if startTimeStr == "" {
		startTimeStr = getOneHour()
	}
	layout := "2006-01-02T15:04"
	t, err := time.Parse(layout, startTimeStr)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	RunInfo(t)
	c.JSON(http.StatusOK, gin.H{
		"message": "Helssssslo, World!",
	})
}

func getOneHour() string {
	currentTime := time.Now()
	oneHourBefore := currentTime.Add(-time.Hour)
	return oneHourBefore.Format("2006-01-02T15:04")
}

func RunInfo(t time.Time) {
	sshConfig := &ssh.ClientConfig{
		User: config.ConfigInstance.FTPConfig.FTPUserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.ConfigInstance.FTPConfig.FTPPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// SSH 连接
	sshClient, err := ssh.Dial("tcp", config.ConfigInstance.FTPConfig.FTPHost+":"+strconv.Itoa(config.ConfigInstance.FTPConfig.FTPPort), sshConfig)
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
	// C:\workspace\Go_download\download
	tmpPath := config.ConfigInstance.LocalDir
	for _, fileInfo := range config.ConfigInstance.CollectorConfig.MeasFileInfo {
		resultList, _ := ListFiles(sftpClient, fileInfo.RemotePath, t)
		fmt.Println("out :", len(resultList))

		// 创建本地文件
		localDirPath := tmpPath + checkOS() + fileInfo.ID
		// 当在本地测试的时候，路径的分隔符是不同的
		//获取本地文件List的方法：
		localFileList, err := GetFileFromFolder(localDirPath)
		if err != nil {
			return
		}
		deltaResulList := SubtractList(resultList, localFileList)
		fmt.Println("======", localDirPath)
		// 检查文件夹是否存在
		_, err = os.Stat(localDirPath)
		if os.IsNotExist(err) {
			// 文件夹不存在，创建文件夹
			err := os.Mkdir(localDirPath, 0755)
			if err != nil {
				fmt.Println("无法创建文件夹:", err)
			}
			fmt.Println("文件夹已创建")
		}
		for _, fileName := range deltaResulList {
			// 读取远端的文件
			remoteFilePath := filepath.Join(fileInfo.RemotePath, fileName)
			fmt.Println("remoteFilePath:", remoteFilePath)
			//仅在window请求linux不耐做出的选择， 在vagrant开发不会发生这种情况
			replacedStr := strings.ReplaceAll(remoteFilePath, "\\", "/")
			remoteFile, err := sftpClient.Open(replacedStr)
			if err != nil {
				log.Fatal(err)
			}
			defer remoteFile.Close()

			localFile, err := os.Create(filepath.Join(localDirPath, fileName))
			if err != nil {
				log.Fatal(err)
			}
			defer localFile.Close()
			// 下载文件内容
			_, err = io.Copy(localFile, remoteFile)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func GetFileFromFolder(folderPath string) ([]string, error) {
	var fileNames []string
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileNames = append(fileNames, file.Name())
	}
	return fileNames, nil
}

func SubtractList(a []string, b []string) []string {
	result := []string{}
	bMap := make(map[string]bool, len(b))
	for _, fileName := range b {
		bMap[fileName] = true
	}
	for _, aFileName := range a {
		if !bMap[aFileName] {
			result = append(result, aFileName)
		}
	}
	return result
}

func ListFiles(sc *sftp.Client, remoteDir string, filtertime time.Time) (resultList []string, err error) {

	files, err := sc.ReadDir(remoteDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to list remote dir: %v\n", err)
		return
	}
	resultList = []string{}
	for _, f := range files {
		var name, modTime, size string

		name = f.Name()
		modTime = f.ModTime().Format("2006-01-02 15:04:05")
		size = fmt.Sprintf("%12d", f.Size())
		fmt.Println(f.Name(), f.ModTime(), filtertime, f.ModTime().Before(filtertime))
		if f.IsDir() || f.ModTime().Before(filtertime) {
			continue
		}
		resultList = append(resultList, name)
		// Output each file name and size in bytes
		fmt.Fprintf(os.Stdout, "%19s %12s %s\n", modTime, size, name)
	}
	// fmt.Println("===", resultList)
	return resultList, nil
}

func checkOS() string {
	switch os := runtime.GOOS; os {
	case "darwin":
		return "/"
	case "linux":
		return "/"
	case "windows":
		return "\\"
	default:
		return "/"
	}
}
