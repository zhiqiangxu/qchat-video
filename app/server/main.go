package main

import "github.com/zhiqiangxu/qchat-video/pkg/api"

func main() {
	server := api.NewServer()
	server.Start()
}
