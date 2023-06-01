package configs

import (
	"fmt"
	"log"
	"testing"
)

func TestReadConfigFile(t *testing.T) {
	config, err := ReadConfigFile("config.yaml")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(config.FTPConfig.FTPHost)
	fmt.Println(config.FTPConfig.FTPUserName)
	fmt.Println(config.FTPConfig.FTPPassword)
	fmt.Println(config.FTPConfig.FTPPort)

	fmt.Println(config.LocalDir)
	for _, fileInfo := range config.CollectorConfig.MeasFileInfo {
		fmt.Println("ID: ", fileInfo.ID)
		fmt.Println("ID: ", fileInfo.FilePattern)
		fmt.Println("ID: ", fileInfo.RemotePath)
	}
}
