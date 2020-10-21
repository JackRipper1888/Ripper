package handler

import (
	"net"
	"time"

	"tools/ctxkit"
	"tools/logkit"
	"tools/mapkit"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/golang/protobuf/proto"

	"Ripper/constant"
	"Ripper/models"
)

var (
	peerList   = mapkit.NewConcurrentSyncMap(64)
	streamInfo = mapkit.NewConcurrentSyncMap(64)

	findNodeResponseChan = make(chan *models.FindNodeResponse, 1)

	peerRequest  = make(chan *models.RequrestInfo, 1)
	peerResponse = make(chan *models.ResponseInfo, 1)
	Conn         *net.UDPConn
)

func MonitorTask() {
	logkit.Info("task|MonitorTask running")
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
func HandleResponseTask() {
	logkit.Info("task|HandleResponseTask running")
	ctx, _ := ctxkit.CtxAdd()
	for {
		select {
		case Val := <-peerRequest:
			//logkit.Succf("peer_ips:%s state %v", Val.Addr.String(), Val.Data[:Val.CountTotal])
			switch Val.Data[0] {
			case constant.HEARTBEAT:
				HeartBeat(Val)
			case constant.FIND_NODE:
				FindNode(Val)
			case constant.FIND_NODE_RESPONSE:
				FindNodeResponse(Val)
			case constant.FIND_VALUE:
				FindValue(Val)
			case constant.FIND_VALUE_RESPONSE:
				FindValueResponse(Val)
			// default:
			// 	body := `{"code": -1,"msg": "cmd not find"}`
			// 	logkit.Err(body)
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
func pushMsg(Addr string, body []byte)  {
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
func FindNodeResponseTask(localId peer.ID) {
	for {
		select {
		case body := <-findNodeResponseChan:
			if rt.Size() < 2000{
				for _, v := range body.Peerlist {
					body, err := proto.Marshal(&models.FindNode{PeerId: string(localId), PeerCount: 7})
					if err != nil {
						logkit.Err(err)
						return
					}
					// data := encodekit.AESEncrypt(body,[]byte("1823eyachlkajsdk"))
					pushMsg(v.Addr,append([]byte{0x25}, body...))
				}
			}
		case now := <-time.After(5* time.Millisecond):
			//push heartbeat
			nowTime := now.Unix()
			timeOut := nowTime - int64(10*time.Millisecond)
			peerList.Range(func (k string,v interface{}) bool {
				peerInfo := v.(*models.PeerInfo)
				if peerInfo.GetTimeStamp() < timeOut{	
					peerList.Delete(k)
					rt.RemovePeer(peer.ID(k))
				}else{
					body, err := proto.Marshal(&models.HeartBeat{PeerId: string(localId), TimeStamp:nowTime })
					if err != nil {
						logkit.Err(err)
						return true
					}
					pushMsg(peerInfo.Addr,append([]byte{0x25}, body...))
				}
				return true
			})
		}
	}
}

