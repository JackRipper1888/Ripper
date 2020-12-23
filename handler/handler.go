package handler

import (
	"context"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"

	"Ripper/constant"
	"Ripper/models"
	"Ripper/providers"
	"Ripper/retrieve"

	"github.com/JackRipper1888/killer/logkit"
)

func HeartBeat(val *models.RequrestInfo) {
	var body models.HeartBeat
	err := proto.Unmarshal(val.Data[1:val.CountTotal], &body)
	if err != nil {
		logkit.Err(err)
		return
	}
	//whitelist checking
	_, isExit := constant.PeerList.Get(string(body.PeerId))
	if isExit {
		_, err = constant.LocalRT.TryAddPeer(peer.ID(body.PeerId), true, true)
		if err != nil {
			logkit.Err(err)
			return
		}
		constant.PeerList.Set(string(body.PeerId), &models.PeerInfo{
			PeerId:    body.PeerId,
			Addr:      val.Addr.String(),
			TimeStamp: body.GetTimeStamp(),
		})
	}
}

func FindNode(val *models.RequrestInfo) {
	var body models.FindNode
	err := proto.Unmarshal(val.Data[1:val.CountTotal], &body)
	if err != nil {
		logkit.Err(err)
		return
	}
	peeridList := constant.LocalRT.NearestPeers([]byte(body.PeerId), int(body.PeerCount))

	peerInfolist := make([]*models.PeerInfo, 0)
	for _, peerId := range peeridList {
		values, isOk := constant.PeerList.Get(string(peerId))
		if !isOk || peerId == peer.ID(body.PeerId) {

			continue
		}
		peerInfolist = append(peerInfolist, (values.(*models.PeerInfo)))
	}

	_, err = constant.LocalRT.TryAddPeer(peer.ID(body.PeerId), true, true)
	if err != nil {
		logkit.Err(err)
		return
	}

	logkit.Succf("FindNode Add peer:%s From peer_ips:%s", string(body.PeerId), val.Addr.String())

	constant.PeerList.Set(string(body.PeerId), &models.PeerInfo{
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
	//logkit.Succ(body.Peerlist, " From peer_addr:", val.Addrç.String())
	//查找任务池发起
	resultList := make([]*models.PeerInfo, 0)
	for _, peerInfo := range body.Peerlist {
		_, err = constant.LocalRT.TryAddPeer(peer.ID(peerInfo.PeerId), true, true)
		if err != nil {
			logkit.Err(err)
			return
		}
		isOk := constant.PeerList.LoadOrStore(string(peerInfo.PeerId), peerInfo)
		if isOk {
			continue
		}
		//同步资源信息（将本地资源 拷贝到距离比自己近的节点上）

		//logkit.Succf("FindNodeResponse Add peer:%s From peer_ips:%s", string(peerInfo.PeerId), val.Addr.String())
		resultList = append(resultList, peerInfo)
	}

	if constant.LocalRT.Size() < constant.TABLE_SIZE {
		localId := retrieve.ConvertPeerID(constant.LocalID)
		for _, v := range resultList {
			body, err := proto.Marshal(&models.FindNode{PeerId: []byte(localId), PeerCount: constant.FIND_NODE_SIZE})
			if err != nil {
				logkit.Err(err)
				return
			}
			//logkit.Succf("FindNode %s request %s ", string(localId), v.Addr)
			pushMsg(v.Addr, append([]byte{constant.FIND_NODE}, body...))
		}
	}

}

func FindProviders(val *models.RequrestInfo) {
	var body models.FindProviders
	err := proto.Unmarshal(val.Data[1:val.CountTotal], &body)
	if err != nil {
		logkit.Err(err)
		return
	}
	body.Leve = body.Leve + 1

	requestChan := make(chan []peer.ID, 1)
	wg := sync.WaitGroup{}
	if body.Leve <= constant.REQUEST_LEVE_LIMIT {
		wg.Add(1)
		go func() {
			defer wg.Done()
			nearpeers := <-requestChan
			for _, peerid := range nearpeers {
				p, isExit := constant.PeerList.Get(string(peerid))
				if isExit {
					requestData, err := proto.Marshal(&body)
					if err != nil {
						logkit.Err(err)
						return
					}
					pushMsg(p.(*models.PeerInfo).GetAddr(), append([]byte{constant.FIND_PROVIDERS}, requestData...))
				}
			}
		}()
	}

	ctx := context.Background()
	peerList := providers.Pm.GetProviders(ctx, []byte(body.Key))
	if len(peerList) == 0 {
		requestChan <- constant.LocalRT.NearestPeers(retrieve.ID(body.Key), constant.FIND_VALUES_SIZE)
	} else {
		requestChan <- peerList

		peerlist := make([]*models.PeerInfo, len(peerList))
		for _, peerid := range peerList {
			p, isExit := constant.PeerList.Get(string(peerid))
			if isExit {
				peerlist = append(peerlist, p.(*models.PeerInfo))
			}
		}
		var responseBody models.FindValueResponse
		responseBody.Peerlist = peerlist
		responseBody.Leve = body.Leve
		responseBody.Key = body.Key
		resultData, err := proto.Marshal(&responseBody)
		if err != nil {
			logkit.Err(err)
			return
		}
		pushMsg(body.GetAddr(), append([]byte{constant.FIND_PROVIDERS_RESPONSE}, resultData...))
	}
	wg.Wait()
}

func FindProvidersResponse(val *models.RequrestInfo) {
	var body models.FindProvidersResponse
	err := proto.Unmarshal(val.Data[1:val.CountTotal], &body)
	if err != nil {
		logkit.Err(err)
		return
	}

	ctx := context.Background()
	for _, v := range body.Peerlist {
		providers.Pm.AddProvider(ctx, []byte(body.Key), peer.ID(v.PeerId))
		streamPeerInfo.Set(string(v.PeerId), v)
	}
	pub.Publish(body)
}

func FindValue(val *models.RequrestInfo) {
	var body models.FindValue
	err := proto.Unmarshal(val.Data[1:val.CountTotal], &body)
	if err != nil {
		logkit.Err(err)
		return
	}

	var findValueResponse models.FindValueResponse
	findValueResponse.Key = body.Key
	findValueResponse.Peerlist = []*models.PeerInfo{}

	//find_store
	findValuesTopic := pub.SubscribeTopic(func(v interface{}) bool {
		if msg, ok := v.(models.FindProvidersResponse); ok {
			return string(msg.Key) == string(body.Key)
		}
		return false
	})
	defer pub.Evict(findValuesTopic)

	wg := sync.WaitGroup{}
	go func() {
		peers := constant.LocalRT.NearestPeers(retrieve.ID(body.Key), constant.FIND_VALUES_SIZE)
		for _, peerid := range peers {
			p, isExit := constant.PeerList.Get(string(peerid))
			if isExit {
				requestData, err := proto.Marshal(&models.FindProviders{
					Key:  body.Key,
					Addr: constant.LISTEN_ADDR,
					Leve: 0,
				})
				if err != nil {
					logkit.Err(err)
					return
				}
				pushMsg(p.(*models.PeerInfo).GetAddr(), append([]byte{constant.FIND_PROVIDERS}, requestData...))
			}
		}
	}()

	wg.Add(1)
	go func ()  {
		defer wg.Done()
		countPeer := 0
		for countPeer >= 10{
			select {
			case data := <-findValuesTopic:
				if msg, ok := data.(models.FindProvidersResponse); ok {
					countPeer += len(msg.Peerlist)
					findValueResponse.Key = msg.Key
					findValueResponse.Leve = msg.Leve
					findValueResponse.Peerlist = msg.Peerlist
				}
			case <-time.After(2 * time.Second):
				break
			}
			findvalueData, err := proto.Marshal(&findValueResponse)
			if err != nil {
				logkit.Err(err)
				return
			}
			pushMsg(val.Addr.String(), append([]byte{constant.FIND_VALUE_RESPONSE}, findvalueData...))
		}
	}()
	wg.Wait()
	return
}

func FindNearUser(val *models.RequrestInfo) {
	var body models.FindNearUser
	err := proto.Unmarshal(val.Data[1:val.CountTotal], &body)
	if err != nil {
		logkit.Err(err)
		return
	}

	peeridList := constant.LocalRT.NearestPeers([]byte(body.PeerId), constant.FIND_VALUES_SIZE)
	if err != nil {
		logkit.Err(err)
		return
	}
	leve := body.Leve + 1

	if len(peeridList[0]) > 0 && peeridList[0] == constant.LocalID {
		peers := make([]*models.PeerInfo, len(peeridList))
		for _, peerid := range peeridList {
			p, isExit := constant.PeerList.Get(string(peerid))
			if isExit {
				peers = append(peers, p.(*models.PeerInfo))
			}
		}
		resultData, err := proto.Marshal(&models.FindNearUserResponse{
			PeerId:   body.PeerId,
			Leve:     leve,
			Peerlist: peers,
		})
		if err != nil {
			logkit.Err(err)
		}
		pushMsg(val.Addr.String(), resultData)
		return
	}

	for _, peerId := range peeridList[:3] {
		values, isOk := constant.PeerList.Get(string(peerId))
		if !isOk || peerId == peer.ID(body.PeerId) {
			continue
		}
		requestData, err := proto.Marshal(&models.FindNearUser{
			PeerId: []byte(peerId),
			Addr:   values.(*models.PeerInfo).GetAddr(),
			Leve:   leve,
		})
		if err != nil {
			logkit.Err(err)
		}
		pushMsg(values.(*models.PeerInfo).GetAddr(), append([]byte{constant.FIND_NEAR_USER}, requestData...))
	}
}

func FindNearUserResponse(val *models.RequrestInfo) {
	var body models.FindNearUserResponse
	err := proto.Unmarshal(val.Data[1:val.CountTotal], &body)
	if err != nil {
		logkit.Err(err)
		return
	}
	pub.Publish(body)
}

func FindUser(val *models.RequrestInfo) {
	var body models.FindUser
	err := proto.Unmarshal(val.Data[1:val.CountTotal], &body)
	if err != nil {
		logkit.Err(err)
		return
	}

	var findValueResponse models.FindUserResponse
	findValueResponse.Key = body.Key
	findValueResponse.Peerlist = []*models.PeerInfo{}

	peeridList := constant.LocalRT.NearestPeers([]byte(body.Key), constant.FIND_VALUES_SIZE)
	if err != nil {
		logkit.Err(err)
		return
	}
	leve := body.Leve + 1

	if len(peeridList) > 0 && peeridList[0] == constant.LocalID {
		peers := make([]*models.PeerInfo, len(peeridList))
		for _, peerid := range peeridList {
			p, isExit := constant.PeerList.Get(string(peerid))
			if isExit {
				peers = append(peers, p.(*models.PeerInfo))
			}
		}
		resultData, err := proto.Marshal(&models.FindNearUserResponse{
			PeerId:   []byte(constant.LocalID),
			Leve:     leve,
			Peerlist: peers,
		})
		if err != nil {
			logkit.Err(err)
		}
		pushMsg(val.Addr.String(), resultData)
		return
	}

	findValuesTopic := pub.SubscribeTopic(func(v interface{}) bool {
		if msg, ok := v.(models.FindNearUserResponse); ok {
			return string(msg.PeerId) == string(body.Key)
		}
		return false
	})
	defer pub.Evict(findValuesTopic)

	go func() {
		for _, peerId := range peeridList[:3] {
			values, isOk := constant.PeerList.Get(string(peerId))
			if !isOk || peerId == peer.ID(body.Key) {
				continue
			}
			requestData, err := proto.Marshal(&models.FindNearUser{
				PeerId: []byte(body.Key),
				Addr:   constant.LISTEN_ADDR,
				Leve:   body.Leve,
			})
			if err != nil {
				logkit.Err(err)
			}
			pushMsg(values.(*models.PeerInfo).GetAddr(),append([]byte{constant.FIND_NEAR_USER},requestData...))
		}
	}()

	wg := sync.WaitGroup{}
	go func() {
		defer wg.Done()
		countPeer := 0
		for countPeer >= 10{
			select {
			case data := <-findValuesTopic:
				if msg, ok := data.(models.FindUserResponse); ok {
					countPeer += len(msg.Peerlist)
					findValueResponse.Key = msg.Key
					findValueResponse.Leve = msg.Leve
					findValueResponse.Peerlist = msg.Peerlist
				}
			case <-time.After(2 * time.Second):
				break
			}
			findvalueData, err := proto.Marshal(&findValueResponse)
			if err != nil {
				logkit.Err(err)
				return
			}
			pushMsg(val.Addr.String(), append([]byte{constant.FIND_USER_RESPONSE}, findvalueData...))
		}
	}()
	
}

func Cache(val *models.RequrestInfo) {
	var body models.FindNearUserResponse
	err := proto.Unmarshal(val.Data[1:val.CountTotal], &body)
	if err != nil {
		logkit.Err(err)
		return
	}
}