package pool

import (
	"fmt"
	config "godownload/configs"
	"io"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var SFTPPool *GenericPool

func LoadPool() {
	connConfig := config.DownloadSFTP.GetConnConfig()
	factory := func() (io.Closer, error) {
		config := &ssh.ClientConfig{
			User: connConfig.FTPUserName,
			Auth: []ssh.AuthMethod{
				ssh.Password(connConfig.FTPPassword),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", connConfig.FTPHost, connConfig.FTPPort), config)
		if err != nil {
			return nil, err
		}
		sftpClient, err := sftp.NewClient(client)
		if err != nil {
			return nil, err
		}
		return sftpClient, nil
	}
	var er error
	SFTPPool, er = NewGenericPool(1, 5, time.Minute, factory)
	if er != nil {
		fmt.Println("err happened in initilize connect pool: ", er)
	}
}
