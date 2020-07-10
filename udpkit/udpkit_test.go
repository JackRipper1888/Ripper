package udpkit

import (
	"Ripper/ctxkit"
	"github.com/astaxie/beego/logs"
	"sync"
	"testing"
)

var wg = sync.WaitGroup{}

func TestListenUdpTask(t *testing.T) {
	//demo()
	//开启udp客户端
	//peerclient()
	//开启udp服务端监听
	go ListenUdpTask()
	//开启udp服务端处理池
	//for i := 0; i<4;i++ {
	//	Worker(i)
	//}
	////开启返回消息协程
	//go RepPeerCmdTask()
	//time.Sleep(1*time.Second)
	wg.Add(1)
	wg.Wait()
	ctxkit.CancelAll()
}

var (
	userInfo = map[string]interface{}{
		"cmd": "userinfo",
		"data": map[string]interface{}{
			"uid":        "123h12kjhjkf",
			"all_stream": "[]",
			"locX":       "10.999999",
			"locY":       "10.999997"}}
	geoNear = map[string]interface{}{
		"cmd": "geo_near",
		"data": map[string]interface{}{
			"point_X": 104.00,
			"point_Y": 20.00,
		}}
	peer_state = map[string]interface{}{
		"cmd":      "state",
		"appid":    "asagdquhw",
		"platform": 4,
		"data": map[string]interface{}{
			"source_recv_len": 6000,
			"p2p_recv_len":    1024,
			"p2p_send_len":    1024,
		},
	}
	peer_state2 = map[string]interface{}{
		"cmd":      "state",
		"appid":    "asagdquhw",
		"platform": 2,
		"data": map[string]interface{}{
			"source_recv_len": 6000,
			"p2p_recv_len":    1024,
			"p2p_send_len":    1024,
		},
	}
	peer_state3 = map[string]interface{}{
		"cmd":      "state",
		"appid":    "666",
		"platform": 4,
		"data": map[string]interface{}{
			"source_recv_len": 6000,
			"p2p_recv_len":    1024,
			"p2p_send_len":    1024,
		},
	}
	tracker_addr = "192.168.1.215:8091"

	serveraddr  = "183.60.143.82:5000"
	serveraddr1 = "192.168.100.200:5000"
)

//go test -v --run TestUdpClient udpkit/udpkit_test.go udpkit/udpkit.go
func TestUdpClient(t *testing.T) {
	i := 0
	for i < 100 {
		_, err := UdpClent(tracker_addr, peer_state)
		if err != nil {
			logs.Error(err)
		}
		//_, err = UdpClent(tracker_addr, peer_state2)
		//if err != nil {
		//	logs.Error(err)
		//}
		//_, err = UdpClent(tracker_addr, peer_state3)
		//if err != nil {
		//	logs.Error(err)
		//}
		i++
	}
	//_, err := UdpClent(tracker_addr, peer_state)
	//if err != nil {
	//	logs.Error(err)
	//}
	//resultInfo := make(map[string]interface{}, 0)
	//json.Unmarshal(data, &resultInfo)
	//fmt.Println("reading", resultInfo)
	//
	////Tcpclient(9000)
	//time.Sleep(24 * time.Hour)
}
