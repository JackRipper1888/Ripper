package udpkit

import (
	"encoding/json"
	"fmt"
	"github.com/Ripper/ctxkit"
	"github.com/astaxie/beego/logs"
	"sync"
	"testing"
	"time"
)

var wg = sync.WaitGroup{}

//
//var (
//
//	login = map[string]interface{}{
//		"cmd":      "login",
//		"lanips": []string{"192.168.100.254:44471"},
//		"nattype":    1,
//	}
//	keepalive = map[string]interface{}{
//		"cmd":      "keepalive",
//		"streamid": "05071246553639577026050703298691",
//		"nattype":    1,
//	}
//	peerlist = map[string]interface{}{
//		"cmd":      "peerlist",
//		"streamid": "05071246553639577026050703298691",
//		"order":    1,
//	}
//	cache = map[string]interface{}{
//		"cmd":      "peerlist",
//		"data": []map[string]interface{}{
//			{
//				"streamid":"05071246553639577026050703298691",
//				"order":1,
//				"lasttime":1590030441,
//			},
//		},
//		"nattype":    1,
//	}
//	body map[string]interface{} = cache
//)
//"183.60.143.82:3030"
//func peerclient() {
//	i := 0
//	wg :=sync.WaitGroup{}
//	for i < 100 {
//		wg.Add(1)
//		i++
//		go func(i int){
//			logs.Info("peerclient() start",i)
//			defer wg.Add(-1)
//			remoteip := net.ParseIP("183.60.143.82")
//			rAddr := &net.UDPAddr{IP: remoteip, Port: 8001}
//
//			conn, err := net.DialUDP("udp", nil, rAddr)
//			if err != nil {
//				fmt.Println(err)
//				return
//			}
//			msgNum := 0
//			for msgNum < 100000{
//				msgNum++
//				body, _ := json.Marshal(body)
//				fmt.Printf("%d号udp消息:%d \n",i,msgNum)
//				if _, err := conn.Write([]byte(body)); err != nil {
//					fmt.Println(err)
//					return
//				}
//
//				msg := make([]byte, 1024)
//				conn.Read(msg)
//
//			}
//			conn.Close()
//		}(i)
//	}
//	wg.Wait()
//}
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
	serveraddr  = "183.60.143.82:5000"
	serveraddr1 = "192.168.100.200:5000"
)

func TestUdpClient(t *testing.T) {
	data, err := UdpClent(serveraddr, geoNear)
	if err != nil {
		logs.Error(err)
	}
	resultInfo := make(map[string]interface{}, 0)
	json.Unmarshal(data, &resultInfo)
	fmt.Println("reading", resultInfo)

	//Tcpclient(9000)
	time.Sleep(24 * time.Hour)
}
