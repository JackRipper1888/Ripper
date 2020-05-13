package udpkit

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"net"
	"os"
	"sync"
	"testing"
)

var wg = sync.WaitGroup{}

//"183.60.143.82:3030"
func peerclient() {
	logs.Info("peerclient() start")
	conn, err := net.Dial("udp", "183.60.143.83:8081")
	if err != nil {
		os.Exit(1)
		return
	}
	defer conn.Close()
	msgNum := 0
	for {
		msgNum++
		data := map[string]interface{}{
			"cmd":      "peerlist",
			"streamid": "hx",
			"order":    2,
		}
		body, _ := json.Marshal(data)
		if _, err := conn.Write([]byte(body)); err != nil {
			fmt.Println(err)
			return
		}
		msg := make([]byte, 1024)
		fmt.Println("reading")
		_, err = conn.Read(msg)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(msg)
	}
}

func TestListenUdpTask(t *testing.T) {
	//demo()
	//开启udp客户端
	peerclient()
	//开启udp服务端监听
	//go ListenUdpTask()
	//开启udp服务端处理池
	//go WokerPool()
	//开启返回消息协程
	//go RepPeerCmdTask()
	//time.Sleep(1*time.Second)
	//ctxkit.CancelAll()
	//	wg.Add(1)
	//	wg.Wait()
}

func TestUdpClent(t *testing.T) {
	//userInfo := map[string]interface{}{
	//	"cmd":"user_info",
	//	"data": map[string]interface{}{
	//		"uid" : "123h12kjhjkf",
	//		"all_stream" : []map[string]interface{}{
	//			{
	//				"stream_id" : "wusha",
	//				"resource_id" : "12g3h1ihd",
	//				"resource_name" : "误杀",
	//			}},
	//		"loc" : [2]float64{10.0, 87.5500030517578}}}
	geoNear := map[string]interface{}{
		"cmd": "geo_near",
		"data": map[string]interface{}{
			"point_X": 10.00,
			"point_Y": 10.00,
		}}
	data, err := UdpClent("192.168.100.200:8091", geoNear)
	if err != nil {
		logs.Error(err)
	}
	resultInfo := make(map[string]interface{}, 0)
	json.Unmarshal(data, &resultInfo)
	fmt.Println("reading", resultInfo)
}
