package main

import (
	"time"
	"context"
	"runtime"

	"Ripper/handler"
	"Ripper/providers"
	"Ripper/retrieve"
	"Ripper/peer"
	"Ripper/constant"
	
	"tools/ctxkit"
	"tools/logkit"
)

func main() {
	logkit.LogInit(constant.LOG_PATH)

	local_id := peer.GenerateId()
	ctx := context.Background()

	go handler.Monitor()
	go handler.ResultCmdTask()

	for i := 0; i < runtime.NumCPU(); i++ {
		go handler.HandleResponse()
	}

	err := handler.InitRoutingTable(ctx, retrieve.ConvertPeerID(local_id))
	if err != nil {
		logkit.Err(err)
		return
	}

	err = providers.InitProvider(ctx, local_id)
	if err != nil {
		logkit.Err(err)
		return
	}
	time.Sleep(20*time.Second)
	ctxkit.CancelAll()
}
