package api

import (
	"fmt"
	config "godownload/configs"
	pool "godownload/pool"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
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
	sftpClientObject, err := pool.SFTPPool.Acquire()
	sftpClient := sftpClientObject.(*sftp.Client)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.SFTPPool.Release(sftpClientObject)

	tmpPath := sftpConfig.GetLocalDir()
	var downInfoList []DownloadInfo
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
		downInfoList = append(downInfoList, resultList...)
	}
	// // single g test
	// for _, job := range downInfoList {
	// 	downloadSingleFile(sftpClient, job)
	// }

	//download with work pool
	runWorkPool(3, len(downInfoList), downInfoList)

}

func runWorkPool(maxWorkers int, numTask int, downInfoList []DownloadInfo) {

	taskChan := make(chan DownloadInfo)
	var wg sync.WaitGroup

	// 启动工作池
	for i := 0; i < maxWorkers; i++ {
		sftpClientObject, _ := pool.SFTPPool.Acquire()
		sftpClient := sftpClientObject.(*sftp.Client)
		go worker(taskChan, &wg, sftpClient)
	}
	// 添加任务到任务通道
	for i := 0; i < numTask; i++ {
		job := downInfoList[i]
		taskChan <- job
		wg.Add(1)
	}

	wg.Wait()
	close(taskChan)
}

// 工作池的工作函数
func worker(jobChan <-chan DownloadInfo, wg *sync.WaitGroup, sftpClient *sftp.Client) {
	for job := range jobChan {
		// 执行文件下载
		downloadSingleFile(sftpClient, job)
		wg.Done()
	}
	defer pool.SFTPPool.Release(sftpClient)
}

func downloadSingleFile(sftpClient *sftp.Client, downloadInfo DownloadInfo) {
	//仅在window请求linux不耐做出的选择， 在vagrant开发不会发生这种情况
	remotePath := downloadInfo.RemotePath
	localPath := downloadInfo.LocalPath
	remoteFile, err := sftpClient.Open(remotePath)
	if err != nil {
		fmt.Println("failed in get remote file err: ", err)
		return
	}
	defer remoteFile.Close()
	localFile, err := os.Create(localPath)
	if err != nil {
		fmt.Println("failed in create local file err: ", err)
		return
	}
	defer localFile.Close()
	// 下载文件内容
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		fmt.Println("failed in download remote file to loacl, err: ", err)
		return
	}
	// unzip the .gz file.
	_, err = os.Stat(localPath)
	if err == nil && strings.HasSuffix(localPath, ".gz") {
		err = decodeGZFile(localPath)
		if err != nil {
			fmt.Println("err in unzip:", err)
			return
		}
	}
}
