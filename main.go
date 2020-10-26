package main

import (
	"flag"
	"context"
	"runtime"

	"Ripper/handler"
	"Ripper/providers"
	"Ripper/retrieve"
	"Ripper/peer"
	"Ripper/constant"
	
	"tools/logkit"
)

func main() {
	flag.Parse()
	logkit.LogInit(constant.LOG_PATH)
	
	local_id := peer.GenerateId()
	
	ctx := context.Background()

	go handler.MonitorTask()

	for i := 0; i < runtime.NumCPU(); i++ {
		go handler.HandleResponseTask()
	}

	handler.ConnListen <- 0
	
	err := handler.InitRoutingTable(ctx, retrieve.ConvertPeerID(local_id), constant.REGISTER_ADDR)
	if err != nil {
		logkit.Err(err)
		return
	}

	go handler.FindNodeResponseTask(retrieve.ConvertPeerID(local_id))

	err = providers.InitProvider(ctx, local_id)
	if err != nil {
		logkit.Err(err)
		return
	}
	
	handler.ResultCmdTask()
}
