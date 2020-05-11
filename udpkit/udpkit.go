package udpkit

import (
	"fmt"
	"github.com/Ripper/ctxkit"
	"github.com/astaxie/beego/logs"
	"net"
	"time"
)

var (
	Conn        *net.UDPConn
	repPeerList = make(chan PeerCmd, 10)
	reqPeerList = make(chan PeerCmd, 10)

	Nowtime = time.Now().String()
)

// peer指令
type PeerCmd struct {
	peerAddr   net.UDPAddr
	countTotal int
	data       []byte
}

// 监听udp端口
func ListenUdpTask() {
	logPre := "|listen|ips=%s|msg:%s"
	listenIps := "127.0.0.1:8080"
	netAddr, err := net.ResolveUDPAddr("udp", listenIps)
	if err != nil {
		logs.Error("@"+logPre, listenIps, err.Error())
		//PressEnterToExit()
	}
	Conn, err = net.ListenUDP("udp", netAddr)
	if err != nil {
		logs.Error("@"+logPre, listenIps, err.Error())
		//PressEnterToExit()
	}
	defer Conn.Close()
	logs.Info("Task:ListenUdpTask() start")
	for {
		// Here must use make and give the lenth of buffer
		data := make([]byte, 1024)
		countTotal, peerAddr, err := Conn.ReadFromUDP(data)
		if err != nil {
			logs.Debug("@"+logPre, listenIps, "ReadFromUDP:"+err.Error())
			continue
		}
		// 插入队列
		var Cmd PeerCmd
		Cmd.peerAddr = *peerAddr
		Cmd.data = data
		Cmd.countTotal = countTotal
		reqPeerList <- Cmd
	}
}

//func WokerPool()  {
//	num := runtime.NumCPU()
//	fmt.Println("CPUNUM:",num)
//	for i:= 0; i<10 ;i++{
//		go	Worker(i)
//	}
//	go addWorder()
//
//}

//func addWorder()  {
//	for{
//
//	}
//}
func Worker(i int) {
	count := 0
	for {
		select {
		case data := <-reqPeerList:
			count++
			fmt.Println(i, string(data.data[:data.countTotal]), count)
			repPeerList <- data
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
		case data := <-repPeerList:
			Conn.WriteToUDP(data.data, &data.peerAddr)
		case <-ctx.Done():
			return
		}
	}
}
