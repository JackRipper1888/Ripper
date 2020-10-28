package main

import (
	"runtime"

	"Ripper/handler"
	"Ripper/providers"
	"Ripper/constant"
	
	"tools/logkit"
	"tools/ctxkit"
)

func main() {
	logkit.LogInit(constant.LOG_PATH)
	
	defer ctxkit.CancelAll()
	go handler.MonitorTask()
	
	for i := 0; i < runtime.NumCPU(); i++ {
		go handler.HandleResponseTask()
	}

	handler.ConnListen <- 0
	
	err := handler.InitRoutingTable()
	if err != nil {
		logkit.Err(err)
		return
	}

	go handler.HeartBeatTask(constant.HEARTBEAT_LATENCY_TIME)

	err = providers.InitProvider()
	if err != nil {
		logkit.Err(err)
		return
	}
	
	handler.ResultCmdTask()
}
