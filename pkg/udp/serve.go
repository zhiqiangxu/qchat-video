package udp

import (
	"fmt"
	"net"
	"sync"
	"time"

	"sync/atomic"

	"github.com/zhiqiangxu/qrpc"

	"github.com/zhiqiangxu/qchat/pkg/instance"
)

// Serve for session
type Serve struct {
	sync.RWMutex
	wg      sync.WaitGroup
	stopped int32
	input   AVStartInput
	ln      *net.UDPConn
}

// NewServe is ctor for Serve
func NewServe(input AVStartInput, ln *net.UDPConn) *Serve {
	return &Serve{input: input, ln: ln}
}

// Start serve
func (sv *Serve) Start() {

	qrpc.GoFunc(&sv.wg, sv.handlePackets)

	return
}

func (sv *Serve) handlePackets() {
	bytes := make([]byte, 65*1024)

	for {
		err := sv.ln.SetReadDeadline(time.Now().Add(time.Second * 2))
		if err != nil {
			instance.Logger().Error("SetReadDeadline err", err)
			continue
		}
		n, remoteAddr, err := sv.ln.ReadFromUDP(bytes)
		if err != nil {
			instance.Logger().Error("ReadFromUDP err", err)
			continue
		}
		if !sv.checkRemoteAddr(remoteAddr) {
			instance.Logger().Error("Got packet from invalid remoteAddr", remoteAddr)
			continue
		}
		sv.parsePacketAndForward(bytes[0:n])

		if atomic.LoadInt32(&sv.stopped) != 0 {
			return
		}
	}

}

func (sv *Serve) checkRemoteAddr(remoteAddr *net.UDPAddr) bool {
	return true
}

func (sv *Serve) parsePacketAndForward(bytes []byte) {

}

// RangeUIDs iterates over all uids in session
func (sv *Serve) RangeUIDs(f func(string)) {
	sv.RLock()

	for _, uid := range sv.input.UIDs {
		f(uid)
	}

	sv.RUnlock()
}

// Stop serve
func (sv *Serve) Stop() (err error) {

	err = sv.ln.Close()
	if err != nil {
		return
	}

	swapped := atomic.CompareAndSwapInt32(&sv.stopped, 0, 1)
	if !swapped {
		sv.Unlock()
		err = fmt.Errorf("already stopped")
		return
	}

	sv.wg.Wait()

	return
}
