package sftpclient

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
)

func UploadFile(client *sftp.Client, localPath, remotePath string) error {
	srcFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := client.Create(remotePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	fmt.Println("Uploaded to", remotePath)
	return nil
}
