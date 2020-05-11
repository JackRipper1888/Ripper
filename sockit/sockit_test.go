package sockit

import (
	"fmt"
	"github.com/Ripper/ctxkit"
	"github.com/astaxie/beego/logs"
	"runtime"
	"testing"
	"time"
)

var (
	reqList = make(chan Peer, 10)
	repList = make(chan Peer, 10)
	udpconn Net
)

func TestNewNet(t *testing.T) {
	go ListenUdpTask()
	go WokerPool()
	go RepPeerCmdTask()
	time.Sleep(1 * time.Minute)
}

// 监听udp端口
func ListenUdpTask() {
	udpconn = NewNet("udp")
	udpconn.Listen("192.168.100.200:8098")
	for {
		select {
		case data := <-udpconn.Read(1024):
			reqList <- data
		case <-time.After(1 * time.Second):
			return
		}
	}
}

func WokerPool() {
	num := runtime.NumCPU()
	fmt.Println("CPUNUM:", num)
	for i := 0; i < 10; i++ {
		go Worker(i)
	}

}
func Worker(i int) {
	count := 0
	for {
		select {
		case data := <-reqList:
			count++
			repList <- data
		case <-time.After(1 * time.Minute):
			return
		}
	}
}

// 返回peer指令
func RepPeerCmdTask() {
	logs.Info("Task:RepPeerCmdTask() start")
	ctx, _ := ctxkit.CtxAdd()
	for {
		select {
		case data := <-repList:
			udpconn.Write(data.Data, "192.168.100.100:8091", data.Addr.String())
		case <-ctx.Done():
			return
		}
	}
}
