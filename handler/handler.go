package handler

import (
	"context"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"

	"tools/ctxkit"
	"tools/logkit"
	"tools/mapkit"

	"Ripper/constant"
	"Ripper/models"
	p "Ripper/providers"
	"Ripper/retrieve"
)

var (
	peerList   = mapkit.NewConcurrentSyncMap(64)
	streamInfo = mapkit.NewConcurrentSyncMap(64)

	peerRequest  = make(chan *models.RequrestInfo, 1)
	peerResponse = make(chan *models.ResponseInfo, 1)
	Conn         *net.UDPConn
)

func Monitor() {
	logkit.Info("task|Monitor running")
	logPre := "|listen|ips=%s|msg:%s"
	listenIps := ""
	netAddr, err := net.ResolveUDPAddr("udp", listenIps)
	if err != nil {
		logkit.Err("@"+logPre, listenIps, err.Error())
		return
	}

	Conn, err = net.ListenUDP("udp", netAddr)
	if err != nil {
		logkit.Err("@"+logPre, listenIps, err.Error())
		return
	}

	defer Conn.Close()
	for {
		var bf models.RequrestInfo
		bf.CountTotal, bf.Addr, err = Conn.ReadFromUDP(bf.Data[:])
		if err != nil {
			logkit.Err("@"+logPre, listenIps, "ReadFromUDP:"+err.Error())
			continue
		}
		peerRequest <- &bf
	}
}

// 处理peer发送的指令
func HandleResponse() {
	logkit.Info("task|HandleResponse running")
	ctx, _ := ctxkit.CtxAdd()
	for {
		select {
		case Val := <-peerRequest:
			logkit.Succf("peer_ips:%s state %v", Val.Addr.String(), Val.Data[:Val.CountTotal])
			switch Val.Data[0] {
			case constant.FIND_NODE:
				FindNode(Val)
			case constant.FIND_NODE_RESPONSE:
				FindNodeResponse(Val)
			case constant.FIND_VALUE:
				FindValue(Val)
			default:
				body := `{"code": -1,"msg": "cmd not find"}`
				logkit.Err(body)
				// 	pushResultChannle(*Val.Addr, strkit.Str2bytes(body))

			}
		case <-ctx.Done():
			return
		}
	}
}

// 返回peer指令
func ResultCmdTask() {
	logkit.Info("task|ResultCmdTask running")
	ctx, _ := ctxkit.CtxAdd()
	//var Val *models.ResultInfo
	for {
		select {
		case Val := <-peerResponse:
			Conn.WriteToUDP(Val.Data, Val.Addr)
		case <-ctx.Done():
			return
		}
	}
}

func FindNode(val *models.RequrestInfo) {
	var body models.FindNode
	err := proto.Unmarshal(val.Data[1:val.CountTotal], &body)
	if err != nil {
		logkit.Err(err)
		return
	}
	peeridList := rt.NearestPeers([]byte(body.PeerId), int(body.PeerCount))

	peerInfolist := make([]*models.PeerInfo, 0)
	for _, peerId := range peeridList {
		values, isOk := peerList.Get(string(peerId))
		if !isOk {
			continue
		}
		peerInfolist = append(peerInfolist, (values.(*models.PeerInfo)))
	}

	_, err = rt.TryAddPeer(peer.ID(body.PeerId), true, true)
	if err != nil {
		logkit.Err(err)
		return
	}

	ResponseBody := &models.FindNodeResponse{
		PeerId:   body.PeerId,
		Peerlist: peerInfolist,
	}

	resultData, err := proto.Marshal(ResponseBody)
	if err != nil {
		logkit.Err(err)
		return
	}

	peerResponse <- &models.ResponseInfo{
		Addr: val.Addr,
		Data: append([]byte{0x28}, resultData...),
	}
}

func FindNodeResponse(val *models.RequrestInfo) {
	var body models.FindNodeResponse
	err := proto.Unmarshal(val.Data[1:val.CountTotal], &body)
	if err != nil {
		logkit.Err(err)
		return
	}
	//查找任务池发起
	for _, peerInfo := range body.Peerlist {
		_, err = rt.TryAddPeer(peer.ID(peerInfo.PeerId), true, true)
		if err != nil {
			logkit.Err(err)
			return
		}
		peerList.Set(peerInfo.PeerId, peerInfo)
	}
}

func FindValue(val *models.RequrestInfo) {
	var body models.FindValue
	err := proto.Unmarshal(val.Data[1:val.CountTotal], &body)
	if err != nil {
		logkit.Err(err)
		return
	}
	var responseBody models.FindValueResponse
	ctx := context.Background()
	peerList := p.Pm.GetProviders(ctx, []byte(body.Key))

	responseBody.Leve = body.Leve
	responseBody.Key = body.Key
	if len(peerList) == 0 {
		nearpeers := rt.NearestPeers(retrieve.ID(body.Key), 10)
		nearpeerlist := []*models.PeerInfo{}
		for _, peerid := range nearpeers {
			nearpeerlist = append(nearpeerlist, &models.PeerInfo{
				PeerId: string(peerid),
			})
		}
		responseBody.Nearpeerlist = nearpeerlist
	} else {
		peerlist := []*models.PeerInfo{}
		for _, peerid := range peerList {
			peerlist = append(peerlist, &models.PeerInfo{
				PeerId: string(peerid),
			})
		}
		responseBody.Peerlist = peerlist
	}
	resultData, err := proto.Marshal(&responseBody)
	if err != nil {
		logkit.Err(err)
		return
	}
	peerResponse <- &models.ResponseInfo{
		Addr: val.Addr,
		Data: append([]byte{0x28}, resultData...),
	}
}

func FindValueResponse(val *models.RequrestInfo) {
	var body models.FindValueResponse
	err := proto.Unmarshal(val.Data[1:val.CountTotal], &body)
	if err != nil {
		logkit.Err(err)
		return
	}
	ctx := context.Background()
	if len(body.Peerlist) > 0 {
		for _, v := range body.Peerlist {
			p.Pm.AddProvider(ctx, []byte(body.Key), peer.ID(v.PeerId))
		}
	}
	if len(body.Nearpeerlist) > 0 && len(p.Pm.GetProviders(ctx,[]byte(body.Key))) < 7{
		for _, v := range body.Nearpeerlist {
			udpaddr, err := net.ResolveUDPAddr("udp", v.Addr)
			if err != nil {
				logkit.Err(err)
				return
			}
			body, err := proto.Marshal(&models.FindValue{Key: v.PeerId})
			if err != nil {
				logkit.Err(err)
				return
			}
			peerResponse <- &models.ResponseInfo{
				Addr: udpaddr,
				Data: append([]byte{0x28}, body...),
			}
		}
	}

}
