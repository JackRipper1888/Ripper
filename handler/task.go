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
	"Ripper/retrieve"
)

var (
	peerList   = mapkit.NewConcurrentSyncMap(64)
	streamInfo = mapkit.NewConcurrentSyncMap(64)

	findNodeResponseChan = make(chan *models.FindNodeResponse, 1)

	peerRequest  = make(chan *models.RequrestInfo, 1)
	peerResponse = make(chan *models.ResponseInfo, 1)

	ConnListen	= make(chan int)
	Conn         *net.UDPConn
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
func FindNodeResponseTask(localId retrieve.ID) {
	for {
		select {
		case body := <-findNodeResponseChan:
			if rt.Size() < 2000{
				for _, v := range body.Peerlist {
					body, err := proto.Marshal(&models.FindNode{PeerId: []byte(localId), PeerCount: 7})
					if err != nil {
						logkit.Err(err)
						return
					}
					// data := encodekit.AESEncrypt(body,[]byte("1823eyachlkajsdk"))
					logkit.Succf("FindNode %s request %s ",string(localId),v.Addr)
					pushMsg(v.Addr,append([]byte{constant.FIND_NODE}, body...))
				}
			}
		case now := <-time.After(5* time.Minute):
			//push heartbeat
			nowTime := now.Unix()
			timeOut := nowTime - int64(10*time.Minute)
		
			peerList.Range(func (k string,v interface{}) bool {
				peerInfo := v.(*models.PeerInfo)
				if peerInfo.GetTimeStamp() < timeOut{	
					peerList.Delete(k)
					rt.RemovePeer(peer.ID(k))
					logkit.Succ("delete peer_id :",peerInfo.PeerId)
				} else {
					body, err := proto.Marshal(&models.HeartBeat{PeerId: []byte(localId), TimeStamp:nowTime })
					if err != nil {
						logkit.Err(err)
						return true
					}
					logkit.Succf("HeartBeat push peer_ips:%s peer_id:%s",peerInfo.Addr,string(localId))
					pushMsg(peerInfo.Addr,append([]byte{constant.HEARTBEAT}, body...))
				}
				return true
			})
		}
	}
}

