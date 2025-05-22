package config

type ServerConfig struct {
	Host       string
	Port       string
	Username   string
	Password   string
	RemotePath string // path to upload/download file
}

var Servers = []ServerConfig{
	{
		Host:       "ip address",
		Port:       "22",
		Username:   "user1",
		Password:   "yourpassword",
		RemotePath: "/remote/path/file1.txt",
	},
	{
		Host:       "ip",
		Port:       "22",
		Username:   "User",
		Password:   "Pass",
		RemotePath: "/home/test",
	},
}
