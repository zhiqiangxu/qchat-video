package udp

// Server for udp stream server
type Server struct {
}

// NewServer is ctor for Server
func NewServer() *Server {
	return &Server{}
}

// Session for av
type Session struct {
}

// AVStart for start
func (s *Server) AVStart(session Session) (err error) {
	return
}

// AVEnd for end
func (s *Server) AVEnd(session Session) (err error) {
	return
}
