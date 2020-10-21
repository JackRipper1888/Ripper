package handler

import (
	"time"
	"context"
	_"net"

	pstore "github.com/libp2p/go-libp2p-peerstore"
	// "github.com/golang/protobuf/proto"

	// "Ripper/models"
	"Ripper/retrieve"
	"tools/logkit"
)

var (
	bucket_size int = 8
	rt          *retrieve.RoutingTable
)

func InitRoutingTable(ctx context.Context, localId retrieve.ID) error {
	ps := pstore.NewMetrics()
	var err error
	rt, err = retrieve.NewRoutingTable(20, localId, time.Minute, ps, 100*time.Minute, nil)
	if err != nil {
		logkit.Err(err)
		return err
	}

	//拿着自己的NODE_ID和公钥去注册中心（拿到入网节点地址和公钥）

	//拿着自己的通信密钥和NODE_ID 与入网节点的公钥进行加密发送给入网节点 进行注册 返回自己的通信密钥

	// 访问已知节点发送FIND_NODE查询自己（用返回节点更新K-bueck）
	// netAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:6666")
	// if err != nil {
	// 	logkit.Err(err)
	// 	return err
	// }

	// body, err := proto.Marshal(&models.FindNode{PeerId: string(localId), PeerCount: 7})
	// if err != nil {
	// 	logkit.Err(err)
	// 	return err
	// }

	// data := encodekit.AESEncrypt(body,[]byte("1823eyachlkajsdk"))
	// peerResponse <- &models.ResponseInfo{
	// 	Addr: netAddr,
	// 	Data: append([]byte{0x25}, body...),
	// }
	return nil
}


// func findNearPeer(udpAddrs []*net.UDPAddr) {
// 	// 请求已知节点返回临近节点 更新 K-bueck
// 	for _, netAddr := range udpAddrs {
// 		peerResponse <- &models.ResponseInfo{
// 			Addr: netAddr,
// 			Data: []byte{},
// 		}
// 		_,err := rt.TryAddPeer(peer.ID(""),true,true)
// 		if err != nil {
// 			logkit.Err(err)
// 			return
// 		}
// 	}
// }
