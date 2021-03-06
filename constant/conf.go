package constant

import (
	"os"
	"path/filepath"
	"github.com/JackRipper1888/killer/confkit"
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
	Server       server `toml:"server"`
	RoutingTable table  `toml:"routing_table"`
}
type server struct {
	Debug                bool   `toml:"debug"`
	ListenAddr           string `toml:"listen_addr"`
	RegisterAddr         string `toml:"register_addr"`
	LogPath              string `toml:"log_path"`
	HeartbeatLatencyTime int64  `toml:"heartbeat_latency_time"`
	RequestLeveLimit     int32  `toml:"request_leve_limit`
}

type table struct {
	FindNodeSize   int32 `toml:"find_node_size"`
	FindValuesSize int   `toml:"find_values_size"`
	BucketSize     int   `toml:"bucket_size"`
	TableSize      int   `toml:"table_size"`
}
