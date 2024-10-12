package server

import (
	"fmt"

	"github.com/fallrising/goku-api/pkg/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	config *config.Config
}

func NewServer(cfg *config.Config, v1Router *gin.RouterGroup) *Server {
	router := gin.Default()
	router.Any("/api/v1/*path", func(c *gin.Context) {
		v1Router.HandleContext(c)
	})

	return &Server{
		router: router,
		config: cfg,
	}
}

func (s *Server) Run() error {
	return s.router.Run(fmt.Sprintf(":%d", s.config.Port))
}
