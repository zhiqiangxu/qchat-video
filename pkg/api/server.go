package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/qchat-video/pkg/config"
)

// Server for http service
type Server struct {
	app *gin.Engine
}

// NewServer creates a Server
func NewServer() *Server {

	app := gin.New()

	app.POST("/avstart", AVStart)
	app.POST("/avend", AVStart)

	s := &Server{app: app}

	return s
}

// Start server
func (s *Server) Start() error {
	return s.app.Run(config.Load().HTTPAddr)
}
