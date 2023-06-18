package api

import (
	"fmt"
	config "godownload/configs"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
)

func DownloadHanlder(c *gin.Context) {
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
	start := time.Now()
	RunDownload(*t, config.DownloadSFTP)
	delta := time.Since(start)

	msg = fmt.Sprintf("spent %d millsecond", delta.Milliseconds())
	c.JSON(returnCode, gin.H{
		"message": msg,
	})
}

func RunDownload(t time.Time, sftpConfig config.SFTPConfig) {
	//prepare for create sftp
	connConfig := sftpConfig.GetConnConfig()
	sshConfig := GetSshConfig(connConfig.FTPUserName, connConfig.FTPPassword)
	host := connConfig.FTPHost
	port := connConfig.FTPPort
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

	tmpPath := sftpConfig.GetLocalDir()
	for _, fileInfo := range sftpConfig.GetDownloadParams() {
		// 创建本地文件
		localDirPath := tmpPath + checkOS() + fileInfo.ID
		err := CreateLocalPath(localDirPath)
		if err != nil {
			fmt.Println("failed in creat folder err:", err.Error())
			return
		}
		// interesting logic to filter files
		resultList, _ := ListFiles(sftpClient, fileInfo, localDirPath, t)
		// downloald files with sync, one by one
		downloadFileList(sftpClient, resultList, fileInfo, localDirPath)

	}
}

func downloadFileList(sftpClient *sftp.Client, resultList []string, fileInfo config.MeasFileInfo, localDirPath string) error {
	for _, fileName := range resultList {
		// 读取远端的文件
		remoteFilePath := filepath.Join(fileInfo.RemotePath, fileName)
		fmt.Println("remoteFilePath:", remoteFilePath)
		//仅在window请求linux不耐做出的选择， 在vagrant开发不会发生这种情况
		replacedStr := strings.ReplaceAll(remoteFilePath, "\\", "/")
		remoteFile, err := sftpClient.Open(replacedStr)
		if err != nil {
			fmt.Println("failed in get remote file err: ", err)
			continue
		}
		defer remoteFile.Close()

		localFilePath := filepath.Join(localDirPath, fileName)
		localFile, err := os.Create(localFilePath)
		if err != nil {
			fmt.Println("failed in create local file err: ", err)
			continue
		}
		defer localFile.Close()
		// 下载文件内容
		_, err = io.Copy(localFile, remoteFile)
		if err != nil {
			fmt.Println("failed in download remote file to loacl, err: ", err)
			continue
		}
		// unzip the .gz file.
		_, err = os.Stat(localFilePath)
		if err == nil && strings.HasSuffix(localFilePath, ".gz") {
			err = decodeGZFile(localFilePath)
			if err != nil {
				fmt.Println("err in unzip:", err)
				continue
			}
		}
	}
	return nil
}
