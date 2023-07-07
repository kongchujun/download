package api

import (
	"compress/gzip"
	"fmt"
	config "godownload/configs"
	"io"
	"log"
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

func GetSshConfig(username, password string) *ssh.ClientConfig {
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
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

func decodeGZFile(localFilePath string) error {
	file, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("unzip file error: ", err)
		return err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		fmt.Println("unzip file error: ", err)
		return err
	}
	defer gzipReader.Close()

	// new file name
	outputFilePath := strings.TrimSuffix(localFilePath, ".gz")
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return err
	}
	// unzip
	_, err = io.Copy(outputFile, gzipReader)
	if err != nil {
		fmt.Println("Error writing to output file:", err)
		return err
	}
	err = os.Remove(localFilePath)
	if err != nil {
		fmt.Println("delete failed: ", err)
	}
	return nil
}

func GetFileFromFolder(folderPath string) ([]string, error) {
	var fileNames []string
	files, err := os.ReadDir(folderPath)
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

func CreateLocalPath(localDirPath string) error {
	_, err := os.Stat(localDirPath)
	if os.IsNotExist(err) {
		// 文件夹不存在，创建文件夹
		err := os.Mkdir(localDirPath, 0755)
		if err != nil {
			fmt.Println("无法创建文件夹:", err)
		}
		fmt.Println("文件夹已创建")
	}
	return nil
}

type DownloadInfo struct {
	RemotePath string
	LocalPath  string
}

func ListFiles(sc *sftp.Client, fileInfo config.MeasFileInfo, localDirPath string, filtertime time.Time) ([]DownloadInfo, error) {
	// take local files from local dir
	localFileList, err := GetFileFromFolder(localDirPath)
	if err != nil {
		return nil, err
	}
	// option model: filter different condition
	rf, err := NewRemoteFile(sc, fileInfo.RemotePath,
		RemoteFileOptionByTime(filtertime),
		RemoteFileOptionByCompare(localFileList),
		RemoteFileOptionByPattern(fileInfo.FilePattern))
	// rf, err := NewRemoteFile(sc, fileInfo.RemotePath,
	// 	RemoteFileOptionByCompare(localFileList))
	if err != nil {
		return nil, err
	}
	resultList := make([]DownloadInfo, 0, len(rf.files))
	for _, file := range rf.files {
		resultList = append(resultList, DownloadInfo{
			RemotePath: filepath.Join(fileInfo.RemotePath, file.Name()),
			LocalPath:  filepath.Join(localDirPath, file.Name()),
		})
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
