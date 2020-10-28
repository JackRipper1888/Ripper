package constant

const (
	HEARTBEAT = 0x22

	FIND_NODE          = 0x25
	FIND_NODE_RESPONSE = 0x26

	FIND_VALUE          = 0x27
	FIND_VALUE_RESPONSE = 0x28

	FIND_KEY           = 0x23
	FIND_KEY_RESPONESE = 0x24

	PEER_INFO_STORE_PATH   = "/conf/peer_info"
)

var (
	LOG_PATH               = ConfData.Server.LogPath
	LISTEN_ADDR            = ConfData.Server.ListenAddr
	REGISTER_ADDR          = ConfData.Server.RegisterAddr
	HEARTBEAT_LATENCY_TIME = ConfData.Server.HeartbeatLatencyTime

	FINDNODESIZE = ConfData.RoutingTable.FindNodeSize
	BUCKETSIZE   = ConfData.RoutingTable.BucketSize
	TABLESIZE    = ConfData.RoutingTable.TableSize
)
