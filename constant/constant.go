package constant

const (
	HEARTBEAT = 0x22

	FIND_NODE          = 0x23
	FIND_NODE_RESPONSE = 0x24

	FIND_PROVIDERS          = 0x25
	FIND_PROVIDERS_RESPONSE = 0x26

	//user_mod
	FIND_VALUE          = 0x27
	FIND_VALUE_RESPONSE = 0x28

	FIND_NEAR_USER          = 0x29
	FIND_NEAR_USER_RESPONSE = 0x30

	FIND_USER          = 0x31
	FIND_USER_RESPONSE = 0x32

	CACHE			   = 0x33
	// FIND_KEY           = 0x29
	// FIND_KEY_RESPONESE = 0x30

	PEER_INFO_STORE_PATH = "/conf/peer_info"
)

const(
	STREAM_TYPE_M3U8 = ".m3u8"	
)

var (
	LOG_PATH               = ConfData.Server.LogPath
	DEBUG                  = ConfData.Server.Debug
	LISTEN_ADDR            = ConfData.Server.ListenAddr
	REGISTER_ADDR          = ConfData.Server.RegisterAddr
	HEARTBEAT_LATENCY_TIME = ConfData.Server.HeartbeatLatencyTime
	REQUEST_LEVE_LIMIT     = ConfData.Server.RequestLeveLimit

	FIND_NODE_SIZE   = ConfData.RoutingTable.FindNodeSize
	FIND_VALUES_SIZE = ConfData.RoutingTable.FindValuesSize
	BUCKET_SIZE      = ConfData.RoutingTable.BucketSize
	TABLE_SIZE       = ConfData.RoutingTable.TableSize
)
