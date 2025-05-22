package sftpclient

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
)

func DownloadFile(client *sftp.Client, remotePath, localDir string) error {
	remoteFile, err := client.Open(remotePath)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	localPath := filepath.Join(localDir, filepath.Base(remotePath))
	localFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		return err
	}

	fmt.Println("Downloaded to", localPath)
	return nil
}
