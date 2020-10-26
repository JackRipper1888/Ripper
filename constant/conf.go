package constant

import (
	"os"
	"path/filepath"
	"tools/confkit"
)

var (
	ConfData = GetConfInfo()
)

func GetConfInfo() AllConfig {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	confDir := dir + "/conf/server.conf"
	confFilePath := confkit.InitCfgFilePath(confDir)
	config := AllConfig{}
	err := confkit.GetConfig(confFilePath, &config)
	if err != nil {
		panic(err)
	}
	return config
}

type AllConfig struct {
	Server server `toml:"server"`
}
type server struct {
	ListenAddr   string `toml:"listen_addr"`
	RegisterAddr string `toml:"register_addr"`
	Loglevel     int    `toml:"loglevel"`
	LogPath      string `toml:"log_path"`
}
