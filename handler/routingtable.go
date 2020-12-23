package handler

import (
	"github.com/golang/protobuf/proto"

	"Ripper/constant"
	"Ripper/models"
	"Ripper/retrieve"

	"github.com/JackRipper1888/killer/logkit"
)

func InitRoutingTable() error {

	if constant.LocalRT.Size() == 0 {
		//拿着自己的NODE_ID和公钥去注册中心（拿到入网节点地址和公钥）
		//拿着自己的通信密钥和NODE_ID 与入网节点的公钥进行加密发送给入网节点 进行注册 返回自己的通信密钥
	
		// 访问已知节点发送FIND_NODE查询自己（用返回节点更新K-bueck）
		body, err := proto.Marshal(&models.FindNode{PeerId: []byte( retrieve.ConvertPeerID(constant.LocalID)), PeerCount: constant.FIND_NODE_SIZE})
		if err != nil {
			logkit.Err(err)
			return err
		}
		pushMsg(constant.REGISTER_ADDR, append([]byte{constant.FIND_NODE}, body...))
	}
	return nil
}
