package handler

import (
	"Ripper/constant"
	"time"
	"context"
	"net"

	pstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/golang/protobuf/proto"

	"Ripper/models"
	"Ripper/retrieve"
	"tools/logkit"
)

var (
	bucket_size int = 8
	rt          *retrieve.RoutingTable
)

func InitRoutingTable(ctx context.Context, localId retrieve.ID,registerAddr string) error {
	ps := pstore.NewMetrics()
	var err error
	logkit.Info("local_id为",localId)
	
	rt, err = retrieve.NewRoutingTable(20, localId, time.Minute, ps, 100*time.Minute, nil)
	if err != nil {
		logkit.Err(err)
		return err
	}
	//拿着自己的NODE_ID和公钥去注册中心（拿到入网节点地址和公钥）

	//拿着自己的通信密钥和NODE_ID 与入网节点的公钥进行加密发送给入网节点 进行注册 返回自己的通信密钥

	// 访问已知节点发送FIND_NODE查询自己（用返回节点更新K-bueck）
	netAddr, err := net.ResolveUDPAddr("udp", registerAddr)
	if err != nil {
		logkit.Err(err)
		return err
	}
	
	body, err := proto.Marshal(&models.FindNode{PeerId: []byte(localId), PeerCount: 7})
	if err != nil {
		logkit.Err(err)
		return err
	}

	//NewRoutingTabledata := encodekit.AESEncrypt(body,[]byte("1823eyachlkajsdk"))
	peerResponse <- &models.ResponseInfo{
		Addr: netAddr,
		Data: append([]byte{constant.FIND_NODE}, body...),
	}
	return nil
}
