package configs

type SFTPConfig interface {
	GetConnConfig() *FTPConfig
	GetDownloadParams() []MeasFileInfo
	GetLocalDir() string
}
