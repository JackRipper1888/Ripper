package main

import (
	"runtime"

	"Ripper/handler"
	"Ripper/providers"
	"Ripper/constant"
	
	"github.com/JackRipper1888/killer/logkit"
	"github.com/JackRipper1888/killer/ctxkit"
)

func main() {
	logkit.LogInit(constant.LOG_PATH)
	
	defer ctxkit.CancelAll()

	providers.InitProvider()
	handler.MakeListenAddr()
	handler.InitRoutingTable()

	go handler.MonitorTask()
	
	for i := 0; i < runtime.NumCPU(); i++ {
		go handler.Worker()
	}

	go handler.ResultTask()

	go handler.HeartBeatTask(constant.HEARTBEAT_LATENCY_TIME)
}