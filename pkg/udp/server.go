package udp

import (
	"context"
	"fmt"
	"sync"

	"github.com/zhiqiangxu/qchat/pkg/core"
	"github.com/zhiqiangxu/qrpc"
)

// Server for udp stream server
type Server struct {
	sync.Mutex
	serves      map[int16]*Serve
	userSession map[string]int16
	wg          sync.WaitGroup
	ctx         context.Context
	cancelFunc  context.CancelFunc
}

// NewServer is ctor for Server
func NewServer() *Server {

	ctx, cancelFunc := context.WithCancel(context.Background())
	return &Server{ctx: ctx, cancelFunc: cancelFunc, serves: make(map[int16]*Serve), userSession: make(map[string]int16)}
}

// SessionType for session type
type SessionType int

const (
	// STAudio for audio
	STAudio SessionType = iota
	// STVideo for video
	STVideo
)

type (
	// AVStartInput for input
	AVStartInput struct {
		App    int
		UIDs   []string
		Type   SessionType
		Port   int16
		Secret string
	}

	// AVStartOutput for output
	AVStartOutput struct {
		core.BaseResp
	}
)

// AVStart for start
func (s *Server) AVStart(input AVStartInput) (r AVStartOutput) {

	serve := NewServe(input)

	s.Lock()
	if _, ok := s.serves[input.Port]; !ok {
		r.SetBase(core.ErrAlreadyInSession, "port already in use")
		s.Unlock()
		return
	}
	for _, uid := range input.UIDs {
		if _, ok := s.userSession[uid]; ok {
			r.SetBase(core.ErrAlreadyInSession, fmt.Sprintf("user already in session:%s", uid))
			s.Unlock()
			return
		}
	}

	s.serves[input.Port] = serve
	for _, uid := range input.UIDs {
		s.userSession[uid] = input.Port
	}
	s.Unlock()

	qrpc.GoFunc(&s.wg, func() {
		serve.Start()
	})
	return
}

type (
	// AVEndInput for input
	AVEndInput struct {
		Port int16
	}

	// AVEndOutput for output
	AVEndOutput struct {
		core.BaseResp
	}
)

// AVEnd for end
func (s *Server) AVEnd(input AVEndInput) (r AVEndOutput) {
	s.Lock()
	if _, ok := s.serves[input.Port]; !ok {
		r.SetBase(core.ErrSessionNotExists, "session not exists")
		s.Unlock()
		return
	}

	serve := s.serves[input.Port]
	delete(s.serves, input.Port)
	serve.RangeUIDs(func(uid string) {
		delete(s.userSession, uid)
	})

	s.Unlock()

	serve.Stop()
	return
}
