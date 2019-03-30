package udp

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/zhiqiangxu/qchat/pkg/instance"

	"github.com/zhiqiangxu/qchat/pkg/core"
)

// Server for udp stream server
type Server struct {
	sync.Mutex
	shutdown    int32
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

var (
	errShutdown = errors.New("already shutdown")
)

// Shutdown server
func (s *Server) Shutdown() (err error) {
	swapped := atomic.CompareAndSwapInt32(&s.shutdown, 0, 1)
	if !swapped {
		err = errShutdown
		return
	}

	s.Lock()

	for _, serve := range s.serves {
		err = serve.Stop()
		if err != nil {
			instance.Logger().Errorln("Stop err in Shutdown", err)
			return
		}
	}

	s.serves = make(map[int16]*Serve)
	s.userSession = make(map[string]int16)
	s.Unlock()

	return
}

// AVStart for start
func (s *Server) AVStart(input AVStartInput) (r AVStartOutput) {

	if atomic.LoadInt32(&s.shutdown) != 0 {
		r.SetBase(core.ErrAPI, errShutdown.Error())
		return
	}

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

	err := serve.Start()
	if err != nil {
		r.SetBase(core.ErrAPI, err.Error())
		return
	}

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
	if atomic.LoadInt32(&s.shutdown) != 0 {
		r.SetBase(core.ErrAPI, errShutdown.Error())
		return
	}

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
