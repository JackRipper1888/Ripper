package protocol

type Block struct {
	Index     int
	Timestamp string
	StreamId  string
	Hash      string
	PrevHash  string
}

type BlockUpdate struct {
	Index     int
	Cmd       int
	Timestamp string
	StreamId  string
	DataPath  string
	Hash      string
}

type PeerUpdateModels struct {
	Cmd       byte
	Index     [4]byte
	TimeStamp [4]byte
	Hash      [8]byte
	StreamId  [8]byte
	Name      [8]byte
	FilerHash [8]byte
}

type Node struct {
	Index     int32
	Timestamp string
	StreamId  string
	Hash      string
	PrevHash  string
	Next      *Node
}
