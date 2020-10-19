package handler

import (
	"context"
	"net"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"

	"Ripper/models"
	"Ripper/retrieve"
	"tools/logkit"
)

var (
	bucket_size int = 8
	rt          *retrieve.RoutingTable
	//ps          peerstore.Peerstore
)

func InitRoutingTable(ctx context.Context, localId retrieve.ID) error {
	ps := pstore.NewMetrics()
	var err error
	rt, err = retrieve.NewRoutingTable(20, localId, time.Minute, ps, 100*time.Minute, nil)
	if err != nil {
		return err
	}

	// 访问已知节点发送FIND_NODE查询自己（用返回节点更新K-bueck）
	// netAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:6666")
	// if err != nil {
	// 	logkit.Err(err)
	// 	return err
	// }

	// findNearPeer([]*net.UDPAddr{netAddr})
	return nil
}

func findNearPeer(udpAddrs []*net.UDPAddr) {
	// 请求已知节点返回临近节点 更新 K-bueck
	for _, netAddr := range udpAddrs {
		peerResponse <- &models.ResponseInfo{
			Addr: netAddr,
			Data: []byte{},
		}
		
		_,err := rt.TryAddPeer(peer.ID(""),true,true)
		if err != nil {
			logkit.Err(err)
			return
		}

	}
}
