package configs

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type FTPConfig struct {
	FTPHost     string `yaml:"ftpHost"`
	FTPUserName string `yaml:"ftpUserName"`
	FTPPassword string `yaml:"ftpPassword"`
	FTPKeyPath  string `yaml:"ftpKeyPath"`
	Passphrase  string `yaml:"passphrase"`
	FTPPort     int    `yaml:"ftpPort"`
	RemotePath  string `yaml:"remotePath"`
}

type MeasFileInfo struct {
	ID          string `yaml:"id"`
	FilePattern string `yaml:"filePattern"`
	RemotePath  string `yaml:"remotePath"`
}

type CollectorConfig struct {
	MeasFileInfo []MeasFileInfo `yaml:"measFileInfo"`
}

type LocalDir struct {
	LocalPath string `yaml:"localPath"`
}

type Config struct {
	FTPConfig       FTPConfig       `yaml:"ftp"`
	LocalDir        string          `yaml:"localPath"`
	CollectorConfig CollectorConfig `yaml:"collector"`
}

func ReadConfigFile(filePath string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

var ConfigInstance *Config

func LoadConfig(filePath string) (*Config, error) {
	var yamlConfig *Config
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		return nil, err
	}
	return yamlConfig, nil
}
