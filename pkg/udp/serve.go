package udp

import (
	"net"
	"sync"
)

// Serve for session
type Serve struct {
	sync.RWMutex
	input AVStartInput
}

// NewServe is ctor for Serve
func NewServe(input AVStartInput) *Serve {
	return &Serve{input: input}
}

// Start serve
func (sv *Serve) Start() {
	addr := net.UDPAddr{
		Port: int(sv.input.Port),
		IP:   net.ParseIP("0.0.0.0"),
	}
	net.ListenUDP("udp", &addr)
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
func (sv *Serve) Stop() {

}
