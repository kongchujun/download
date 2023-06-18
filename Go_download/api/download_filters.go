package api

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"time"

	"github.com/pkg/sftp"
)

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

func RemoteFileOptionByPattern(pattern string) RemoteFileOption {
	return func(rf *RemoteFile) {
		tmpFiles := make([]fs.FileInfo, 0, len(rf.files))
		counter := 0
		for _, file := range rf.files {
			match, err := regexp.MatchString(pattern, file.Name())
			if err != nil {
				fmt.Println("matching error: ", err.Error())
				continue
			}
			if match {
				tmpFiles = append(tmpFiles, file)
				counter++
			}
		}
		rf.files = tmpFiles
	}
}
