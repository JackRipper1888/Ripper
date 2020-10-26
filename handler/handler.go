package handler

import (
	"context"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"

	"Ripper/constant"
	"Ripper/models"
	"Ripper/providers"
	"Ripper/retrieve"
	"tools/logkit"
)

func HeartBeat(val *models.RequrestInfo) {
	var body models.HeartBeat
	err := proto.Unmarshal(val.Data[1:val.CountTotal], &body)
	if err != nil {
		logkit.Err(err)
		return
	}

	_, err = rt.TryAddPeer(peer.ID(body.PeerId), true, true)
	if err != nil {
		logkit.Err(err)
		return
	}

	logkit.Succf("HeartBeat Set peer:%s", string(body.PeerId))
	peerList.Set(string(body.PeerId), &models.PeerInfo{
		PeerId:    body.PeerId,
		Addr:      val.Addr.String(),
		TimeStamp: body.GetTimeStamp(),
	})
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
		if !isOk || peerId == peer.ID(body.PeerId) {
			logkit.Err(body.PeerId, peerId)
			continue
		}
		//logkit.Err(,peerId,peer.ID(body.PeerId))
		peerInfolist = append(peerInfolist, (values.(*models.PeerInfo)))
	}

	_, err = rt.TryAddPeer(peer.ID(body.PeerId), true, true)
	if err != nil {
		logkit.Err(err)
		return
	}

	logkit.Succf("FindNode Add peer:%s From peer_ips:%s", string(body.PeerId), val.Addr.String())

	peerList.Set(string(body.PeerId), &models.PeerInfo{
		PeerId:    body.PeerId,
		Addr:      val.Addr.String(),
		TimeStamp: time.Now().Unix(),
	})

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
		Data: append([]byte{constant.FIND_NODE_RESPONSE}, resultData...),
	}
}

func FindNodeResponse(val *models.RequrestInfo) {
	var body models.FindNodeResponse
	err := proto.Unmarshal(val.Data[1:val.CountTotal], &body)
	if err != nil {
		logkit.Err(err)
		return
	}
	logkit.Succ(body.Peerlist," From peer_addr:",val.Addr.String())
	//查找任务池发起
	for _, peerInfo := range body.Peerlist {
		_, err = rt.TryAddPeer(peer.ID(peerInfo.PeerId), true, true)
		if err != nil {
			logkit.Err(err)
			return
		}

		logkit.Succf("FindNodeResponse Add peer:%s From peer_ips:%s", string(peerInfo.PeerId), val.Addr.String())
		peerList.Set(string(peerInfo.PeerId), peerInfo)
	}
	findNodeResponseChan <- &body
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
	peerList := providers.Pm.GetProviders(ctx, []byte(body.Key))

	responseBody.Leve = body.Leve + 1
	responseBody.Key = body.Key
	if len(peerList) == 0 {
		nearpeers := rt.NearestPeers(retrieve.ID(body.Key), 10)
		nearpeerlist := []*models.PeerInfo{}
		for _, peerid := range nearpeers {
			nearpeerlist = append(nearpeerlist, &models.PeerInfo{
				PeerId: []byte(peerid),
			})
		}
		responseBody.Nearpeerlist = nearpeerlist
	} else {
		peerlist := []*models.PeerInfo{}
		for _, peerid := range peerList {
			peerlist = append(peerlist, &models.PeerInfo{
				PeerId: []byte(peerid),
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
			providers.Pm.AddProvider(ctx, []byte(body.Key), peer.ID(v.PeerId))
		}
	}
	if body.Leve >= 7 {
		return
	}
	if len(body.Nearpeerlist) > 0 && len(providers.Pm.GetProviders(ctx, []byte(body.Key))) < 7 {
		for _, v := range body.Nearpeerlist {
			udpaddr, err := net.ResolveUDPAddr("udp", v.Addr)
			if err != nil {
				logkit.Err(err)
				return
			}
			body, err := proto.Marshal(&models.FindValue{Key: v.PeerId, Leve: body.Leve})
			if err != nil {
				logkit.Err(err)
				return
			}
			peerResponse <- &models.ResponseInfo{
				Addr: udpaddr,
				Data: append([]byte{0x27}, body...),
			}
		}
	}

}
