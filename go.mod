module Ripper

go 1.14

require (
	github.com/golang/protobuf v1.4.3
	github.com/google/protobuf v3.13.0+incompatible // indirect
	github.com/ipfs/go-datastore v0.4.5
	github.com/libp2p/go-libp2p-core v0.7.0
	github.com/libp2p/go-libp2p-kad-dht v0.10.0
	github.com/libp2p/go-libp2p-kbucket v0.4.7
	github.com/libp2p/go-libp2p-peerstore v0.2.6
	github.com/minio/sha256-simd v0.1.1
	golang.org/x/net v0.0.0-20200923182212-328152dc79b1
	google.golang.org/protobuf v1.25.0
	tools v0.1.4
)

replace tools v0.1.4 => ../tools
