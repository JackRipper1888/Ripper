package handler

import (
	//"syscall"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/JackRipper1888/killer/ctxkit"
	"github.com/JackRipper1888/killer/iokit"
	"github.com/JackRipper1888/killer/logkit"
	"github.com/JackRipper1888/killer/mapkit"

	"github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"

	"Ripper/constant"
	"Ripper/models"
	"Ripper/pubsub"
	"Ripper/retrieve"
)

var (
	pub            = pubsub.NewPublisher(100*time.Millisecond, 10)
	streamPeerInfo = mapkit.NewConcurrentSyncMap(64)

	findNodeResponseChan = make(chan []*models.PeerInfo, runtime.NumCPU())

	peerRequest  = make(chan *models.RequrestInfo, runtime.NumCPU())
	peerResponse = make(chan *models.ResponseInfo, runtime.NumCPU())

	ConnListen = make(chan int)
	Conn       *net.UDPConn
)

func MakeListenAddr() {
	logPre := "|listen|ips=%s|msg:%s"
	listenAddr := constant.LISTEN_ADDR
	if constant.DEBUG {
		//本地绑
		listenAddr = "0.0.0.0:0"
	}
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
}
func MonitorTask() {
	//syscall.ForkExec()
	logkit.Info("task|MonitorTask running")

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	for {
		var bf models.RequrestInfo
		bf.CountTotal, bf.Addr, err = Conn.ReadFromUDP(bf.Data[:])
		if err != nil {
			logkit.Err(err)
			continue
		}
		peerRequest <- &bf
	}
}

// 处理peer发送的指令
func Worker() {
	logkit.Info("task|HandleResponseTask running")
	ctx, _ := ctxkit.CtxAdd()
	for {
		select {
		case Val := <-peerRequest:
			switch Val.Data[0] {
			//节点之间通信
			case constant.HEARTBEAT:
				HeartBeat(Val)
			case constant.FIND_NODE:
				FindNode(Val)
			case constant.FIND_NODE_RESPONSE:
				FindNodeResponse(Val)
			case constant.FIND_PROVIDERS:
				FindProviders(Val)
			case constant.FIND_PROVIDERS_RESPONSE:
				FindProvidersResponse(Val)

			case constant.FIND_NEAR_USER:
				FindNearUser(Val)
			case constant.FIND_NEAR_USER_RESPONSE:
				FindNearUserResponse(Val)

			//外部调用
			case constant.FIND_VALUE:
				FindValue(Val)
			case constant.FIND_USER:
				FindUser(Val)
			case constant.CACHE:
				FindUser(Val)
			}
			runtime.Gosched()
		case <-ctx.Done():
			return
		}
	}
}

// 返回peer指令
func ResultTask() {
	logkit.Info("task|ResultCmdTask running")
	ctx, _ := ctxkit.CtxAdd()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

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
	period := time.Duration(latency_time) * time.Second
	for {
		select {
		case now := <-time.After(period):
			//push heartbeat
			nowTime := now.Unix()
			timeOut := nowTime - 2*latency_time
			peers := make([]*models.PeerInfo, 0, constant.TABLE_SIZE)
			constant.PeerList.Range(func(k string, v interface{}) bool {
				peerInfo := v.(*models.PeerInfo)
				if peerInfo.GetTimeStamp() < timeOut {
					//logkit.Err(nowTime,peerInfo.GetTimeStamp(),timeOut)
					constant.PeerList.Delete(k)
					constant.LocalRT.RemovePeer(peer.ID(k))
					logkit.Succf("delete peer_ips:%s peer_id:%x", peerInfo.Addr, peerInfo.PeerId)
				} else {
					peers = append(peers, peerInfo)
					body, err := proto.Marshal(&models.HeartBeat{PeerId: []byte(localId), TimeStamp: nowTime})
					if err != nil {
						logkit.Err(err)
						return true
					}
					logkit.Succf("HeartBeat push peer_ips:%s peer_id:%x", peerInfo.Addr, peerInfo.PeerId)
					pushMsg(peerInfo.Addr, append([]byte{constant.HEARTBEAT}, body...))
				}
				return true
			})

			//持久化存储
			if len(peers) > 0 {
				dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
				confDir := dir + constant.PEER_INFO_STORE_PATH
				body, err := proto.Marshal(&models.FindNodeResponse{
					PeerId:   []byte(constant.LocalID),
					Peerlist: peers,
				})
				if err != nil {
					logkit.Err(err)
					return
				}

				err = iokit.Write(confDir, body)
				if err != nil {
					logkit.Err(err)
					return
				}
			}

			//insert new peer when the table is missing peer
			if constant.LocalRT.Size() < constant.TABLE_SIZE {
				peerIDs := constant.LocalRT.NearestPeers(localId, int(constant.FIND_NODE_SIZE))
				body, err := proto.Marshal(&models.FindNode{PeerId: []byte(localId), PeerCount: constant.FIND_NODE_SIZE})
				if err != nil {
					logkit.Err(err)
				}
				for _, peerId := range peerIDs {
					peerInfo, isExist := constant.PeerList.Get(string(peerId))
					if isExist {
						pushMsg(peerInfo.(*models.PeerInfo).GetAddr(), append([]byte{constant.FIND_NODE}, body...))
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
