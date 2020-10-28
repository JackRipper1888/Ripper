package handler

import (
	"os"
	"net"
	"time"
	"path/filepath"

	"tools/ctxkit"
	"tools/logkit"
	"tools/mapkit"
	"tools/iokit"

	"github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"

	"Ripper/constant"
	"Ripper/models"
	"Ripper/retrieve"
)

var (

	streamInfo = mapkit.NewConcurrentSyncMap(64)

	findNodeResponseChan = make(chan []*models.PeerInfo, 1)

	peerRequest  = make(chan *models.RequrestInfo, 1)
	peerResponse = make(chan *models.ResponseInfo, 1)

	ConnListen = make(chan int)
	Conn       *net.UDPConn
)

func MonitorTask() {
	logkit.Info("task|MonitorTask running")
	logPre := "|listen|ips=%s|msg:%s"
	listenAddr := constant.LISTEN_ADDR

	netAddr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		logkit.Err("@"+logPre, listenAddr, err.Error())
		return
	}

	Conn, err = net.ListenUDP("udp", netAddr)
	if err != nil {
		logkit.Err("@"+logPre, listenAddr, err.Error())
		return
	}
	<-ConnListen
	defer Conn.Close()
	for {
		var bf models.RequrestInfo
		bf.CountTotal, bf.Addr, err = Conn.ReadFromUDP(bf.Data[:])
		if err != nil {
			logkit.Err("@"+logPre, listenAddr, "ReadFromUDP:"+err.Error())
			continue
		}
		peerRequest <- &bf
	}
}

// 处理peer发送的指令
func HandleResponseTask() {
	logkit.Info("task|HandleResponseTask running")
	ctx, _ := ctxkit.CtxAdd()
	for {
		select {
		case Val := <-peerRequest:
			switch Val.Data[0] {
			case constant.HEARTBEAT:
				//logkit.Succf("HEARTBEAT peer_ips:%s state %v", Val.Addr.String(), Val.Data[0:Val.CountTotal])
				HeartBeat(Val)
			case constant.FIND_NODE:
				//logkit.Succf("FIND_NODE peer_ips:%s state %v", Val.Addr.String(), Val.Data[:Val.CountTotal])
				FindNode(Val)
			case constant.FIND_NODE_RESPONSE:
				//logkit.Succf("FIND_NODE_RESPONSE peer_ips:%s state %v", Val.Addr.String(), Val.Data[:Val.CountTotal])
				FindNodeResponse(Val)
			case constant.FIND_VALUE:
				FindValue(Val)
			case constant.FIND_VALUE_RESPONSE:
				FindValueResponse(Val)
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
func pushMsg(Addr string, body []byte) {
	netAddr, err := net.ResolveUDPAddr("udp", Addr)
	if err != nil {
		logkit.Err(err)
		return
	}
	peerResponse <- &models.ResponseInfo{
		Addr: netAddr,
		Data: body,
	}
}

func HeartBeatTask(latency_time int64) {
	localId := retrieve.ConvertPeerID(constant.LocalID)
	ctx, _ := ctxkit.CtxAdd()
	period := time.Duration(latency_time)*time.Second
	for {
		select {
		case now := <-time.After(period):
			//push heartbeat
			nowTime := now.Unix()
			timeOut := nowTime - 2*latency_time
			peers := make([]*models.PeerInfo,0,constant.TABLESIZE)
			constant.PeerList.Range(func(k string, v interface{}) bool {
				peerInfo := v.(*models.PeerInfo)
				if peerInfo.GetTimeStamp() < timeOut {
					//logkit.Err(nowTime,peerInfo.GetTimeStamp(),timeOut)
					constant.PeerList.Delete(k)
					constant.LocalRT.RemovePeer(peer.ID(k))
					logkit.Succ("delete peer_id :", peerInfo.PeerId)
				} else {
					peers = append(peers,peerInfo)
					body, err := proto.Marshal(&models.HeartBeat{PeerId: []byte(localId), TimeStamp: nowTime})
					if err != nil {
						logkit.Err(err)
						return true
					}
					logkit.Succf("HeartBeat push peer_ips:%s peer_id:%s", peerInfo.Addr, string(localId))
					pushMsg(peerInfo.Addr, append([]byte{constant.HEARTBEAT}, body...))
				}
				return true
			})

			//持久化存储
			if len(peers) > 0{
				dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
				confDir := dir + constant.PEER_INFO_STORE_PATH
				body,err := proto.Marshal(&models.FindNodeResponse{
					PeerId: []byte(constant.LocalID),
					Peerlist: peers,
				})
				if err != nil {
					logkit.Err(err)
					return
				}
	
				err = iokit.Write(confDir,body)
				if err != nil {
					logkit.Err(err)
					return
				}
			}

			//insert new peer when the table is missing peer
			if constant.LocalRT.Size() < constant.TABLESIZE {
				ID := constant.LocalRT.NearestPeers(localId,int(constant.FINDNODESIZE))
				body, err := proto.Marshal(&models.FindNode{PeerId: []byte(localId), PeerCount: constant.FINDNODESIZE})
				if err != nil {
					logkit.Err(err)
				}
				for _, peerId := range ID {
					peerInfo,isExist := constant.PeerList.Get(string(peerId))
					if isExist {
						pushMsg(peerInfo.(*models.PeerInfo).GetAddr(),append([]byte{constant.FIND_NODE}, body...))
					}
				}
			}
		case <- ctx.Done():
			return
		}
	}
}
