package instance

import (
	"sync"

	"github.com/zhiqiangxu/qchat-video/pkg/udp"
)

var (
	instanceUDPServer *udp.Server
	lockUDPServer     sync.Mutex
)

// UDPServer is singleton for udp.Server
func UDPServer() *udp.Server {
	if instanceUDPServer != nil {
		return instanceUDPServer
	}

	lockUDPServer.Lock()
	defer lockUDPServer.Unlock()

	if instanceUDPServer != nil {
		return instanceUDPServer
	}

	instanceUDPServer = udp.NewServer()
	return instanceUDPServer

}
