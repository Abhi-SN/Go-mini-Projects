package main

import (
	config "GoSFTP/models"
	"GoSFTP/sftpclient"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for _, server := range config.Servers {
		wg.Add(1)

		go func(s config.ServerConfig) {
			defer wg.Done()

			client, err := sftpclient.Connect(s.Username, s.Password, s.Host, s.Port)
			if err != nil {
				println("Connection error to", s.Host, ":", err.Error())
				return
			}
			defer client.Close()

			// Upload
			err = sftpclient.UploadFile(client, "./localfile.txt", s.RemotePath)
			if err != nil {
				println("Upload error to", s.Host, ":", err.Error())
			}

			// Download
			err = sftpclient.DownloadFile(client, s.RemotePath, "config.json")
			if err != nil {
				println("Download error from", s.Host, ":", err.Error())
			}
		}(server)
	}

	wg.Wait()
	println("All transfers done.")
}
