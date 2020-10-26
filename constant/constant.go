package constant

const (
	HEARTBEAT = 0x22

	FIND_NODE          = 0x25
	FIND_NODE_RESPONSE = 0x26

	FIND_VALUE          = 0x27
	FIND_VALUE_RESPONSE = 0x28

	FIND_KEY           = 0x23
	FIND_KEY_RESPONESE = 0x24

	TimeFormat = "2006-01-02 15:04:05"
)

var (
	LOG_PATH      = ConfData.Server.LogPath
	LISTEN_ADDR   = ConfData.Server.ListenAddr
	REGISTER_ADDR = ConfData.Server.RegisterAddr
)
