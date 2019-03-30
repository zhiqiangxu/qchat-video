package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/zhiqiangxu/qchat-video/pkg/api"
	"github.com/zhiqiangxu/qchat/pkg/instance"
	"github.com/zhiqiangxu/qrpc"
)

func main() {
	server := api.NewServer()

	ctx, cancelFunc := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	qrpc.GoFunc(&wg, func() {
		err := server.Start()
		cancelFunc()
		if err != nil {
			instance.Logger().Errorln("Start err", err)
		}
	})

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	select {
	case <-quitChan:
		err := server.Stop()
		if err != nil {
			instance.Logger().Errorln("Stop err", err)
		}
	case <-ctx.Done():
		instance.Logger().Infoln("Server stopped")
	}
}
