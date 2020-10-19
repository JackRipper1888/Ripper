package constant

const (
	FIND_NODE           = 0x25
	FIND_NODE_RESPONSE  = 0x26

	FIND_VALUE          = 0x27
	FIND_VALUE_RESPONSE = 0x28

	TimeFormat          = "2006-01-02 15:04:05"
)

var (
	LOG_PATH         = ConfData.Server.LogPath
	TRACKER_LOG_PATH = ConfData.Server.TrackerLogPath
)
