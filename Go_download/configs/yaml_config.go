package configs

import "fmt"

type YamlConfig struct {
	config *Config
}

// check whether realize the impl
var _ SFTPConfig = &YamlConfig{}

var DownloadSFTP *YamlConfig

func NewYamlConfig(path string) error {

	fileConfig, err := LoadConfig(path)
	if err != nil {
		fmt.Println("error in decode yaml config file:", err)
		return nil
	}
	yamlConfig := YamlConfig{
		config: fileConfig,
	}
	DownloadSFTP = &yamlConfig
	return nil
}

func (c *YamlConfig) GetConnConfig() *FTPConfig {
	ftpConfig := c.config.FTPConfig
	return &FTPConfig{
		FTPHost:     ftpConfig.FTPHost,
		FTPUserName: ftpConfig.FTPUserName,
		FTPPassword: ftpConfig.FTPPassword,
		FTPKeyPath:  ftpConfig.FTPKeyPath,
		Passphrase:  ftpConfig.Passphrase,
		FTPPort:     ftpConfig.FTPPort,
		RemotePath:  ftpConfig.RemotePath,
	}
}

func (c *YamlConfig) GetDownloadParams() []MeasFileInfo {
	return c.config.CollectorConfig.MeasFileInfo
}

func (c *YamlConfig) GetLocalDir() string {
	return c.config.LocalDir
}
