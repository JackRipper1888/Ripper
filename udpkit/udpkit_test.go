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

func peerclient() {
	logs.Info("peerclient() start")
	conn, err := net.Dial("udp", "183.60.143.82:8001")
	if err != nil {
		os.Exit(1)
		return
	}
	defer conn.Close()
	msgNum := 0
	for msgNum < 10000 {
		msgNum++
		//select {
		//case <- time.After(1*time.Second):
		data := map[string]interface{}{
			"cmd":  "login",
			"data": "12371892",
		}
		body, _ := json.Marshal(data)
		conn.Write([]byte(body))

		fmt.Println(fmt.Sprintf("Write id:%d...", msgNum))
		//var msg [1024]byte
		//_,err :=conn.Read(msg[0:])
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//fmt.Println("msgcount:",msgNum,string(msg[:]))
	}
}

//"183.60.143.82:3030"
func demo() {
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

func TestUdpTask(t *testing.T) {
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
