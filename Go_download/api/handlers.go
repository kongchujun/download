package api

import (
	"fmt"
	config "godownload/configs"
	"io"
	"io/fs"
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
	// these three value should be taken out from setting in the furture
	// config should be a interface: yaml, postgres, mongo etc.
	timeZone := "Asia/Shanghai"
	queryKey := "start"
	layout := "2006-01-02T15:04"

	//it should be merge into a struct
	returnCode := http.StatusOK
	msg := "ok"

	t, err := GetStartTime(c, queryKey, timeZone, layout)
	if err != nil {
		returnCode = http.StatusBadRequest
		c.JSON(returnCode, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	RunDownload(*t)
	c.JSON(returnCode, gin.H{
		"message": msg,
	})
}

// GetStartTime: get the start time so that program can take how many files you can get
// from remote server.
// todo: untest
// c, _ := gin.CreateTestContext(httptest.NewRecorder())
// 设置请求的 HTTP 方法、路径和查询参数
// c.Request, _ = http.NewRequest("GET", "/user", nil)
// c.Request.URL.RawQuery = "name=John"
func GetStartTime(c *gin.Context, key string, timeZone string, layout string) (*time.Time, error) {
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		return nil, err
	}
	currentTime := time.Now().In(location)
	oneHourBefore := currentTime.Add(-time.Hour)

	startTimeStr := c.DefaultQuery(key, oneHourBefore.Format(layout))

	t, err := time.ParseInLocation(layout, startTimeStr, location)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func GetSshConfig(conf *config.Config) *ssh.ClientConfig {
	sshConfig := &ssh.ClientConfig{
		User: conf.FTPConfig.FTPUserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(conf.FTPConfig.FTPPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return sshConfig
}

func GetSshClient(sshConfig *ssh.ClientConfig, host string, port int) (*ssh.Client, error) {
	sshClient, err := ssh.Dial("tcp", host+":"+strconv.Itoa(port), sshConfig)
	if err != nil {
		log.Fatal(err)
	}
	return sshClient, err
}

func RunDownload(t time.Time) {
	//prepare for create sftp
	sshConfig := GetSshConfig(config.ConfigInstance)
	host := config.ConfigInstance.FTPConfig.FTPHost
	port := config.ConfigInstance.FTPConfig.FTPPort
	sshClient, err := GetSshClient(sshConfig, host, port)
	if err != nil {
		log.Fatal(err)
	}
	defer sshClient.Close()
	// create SFTP client
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		log.Fatal(err)
	}
	defer sftpClient.Close()

	tmpPath := config.ConfigInstance.LocalDir
	for _, fileInfo := range config.ConfigInstance.CollectorConfig.MeasFileInfo {
		// 创建本地文件
		localDirPath := tmpPath + checkOS() + fileInfo.ID
		resultList, _ := ListFiles(sftpClient, fileInfo.RemotePath, localDirPath, t)

		for _, fileName := range resultList {

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
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".error") {
			fileName = strings.Replace(fileName, ".error", "", -1)
		} else if strings.HasSuffix(fileName, ".done") {
			fileName = strings.Replace(fileName, ".done", "", -1)
		}

		fileNames = append(fileNames, fileName)
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

type RemoteFile struct {
	files []fs.FileInfo
}

type RemoteFileOption func(rf *RemoteFile)

func NewRemoteFile(sc *sftp.Client, remoteDir string, opts ...RemoteFileOption) (*RemoteFile, error) {
	files, err := sc.ReadDir(remoteDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to list remote dir: %v\n", err)
		return nil, err
	}
	noDirFile := make([]fs.FileInfo, 0, len(files))
	count := 0
	for _, file := range files {
		if len(file.Name()) == 0 {
			continue
		}
		if !file.IsDir() {
			noDirFile = append(noDirFile, file)
			count++
		}
	}

	res := &RemoteFile{
		files: noDirFile[:count],
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func RemoteFileOptionByTime(filtertime time.Time) RemoteFileOption {
	return func(rf *RemoteFile) {
		tmpFiles := make([]fs.FileInfo, 0, len(rf.files))
		counter := 0
		for _, file := range rf.files {
			if file.ModTime().After(filtertime) {
				tmpFiles = append(tmpFiles, file)
				counter++
			}
		}

		rf.files = tmpFiles[:counter]
	}
}

func RemoteFileOptionByCompare(localFileNames []string) RemoteFileOption {
	return func(rf *RemoteFile) {
		result := make([]fs.FileInfo, 0, len(rf.files))
		bMap := make(map[string]bool, len(localFileNames))
		for _, fileName := range localFileNames {
			bMap[fileName] = true
		}
		count := 0
		for _, file := range rf.files {
			if !bMap[file.Name()] {
				result = append(result, file)
				count++
			}
		}
		rf.files = result[:count]
	}
}

func ListFiles(sc *sftp.Client, remoteDir string, localDirPath string, filtertime time.Time) ([]string, error) {
	// take data from local folder
	_, err := os.Stat(localDirPath)
	if os.IsNotExist(err) {
		// 文件夹不存在，创建文件夹
		err := os.Mkdir(localDirPath, 0755)
		if err != nil {
			fmt.Println("无法创建文件夹:", err)
		}
		fmt.Println("文件夹已创建")
	}
	// take local files from local dir
	localFileList, err := GetFileFromFolder(localDirPath)
	if err != nil {
		return nil, err
	}
	rf, err := NewRemoteFile(sc, remoteDir,
		RemoteFileOptionByTime(filtertime),
		RemoteFileOptionByCompare(localFileList))
	if err != nil {
		return nil, err
	}
	resultList := make([]string, 0, len(rf.files))
	for _, file := range rf.files {
		resultList = append(resultList, file.Name())
	}
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
