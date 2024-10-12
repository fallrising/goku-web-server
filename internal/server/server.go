package server

import (
	"fmt"

	"github.com/fallrising/goku-api/internal/database"
	"github.com/fallrising/goku-api/internal/handlers"
	"github.com/fallrising/goku-api/pkg/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	config *config.Config
	db     *database.Database
}

func NewServer(cfg *config.Config, db *database.Database) *Server {
	return &Server{
		router: gin.Default(),
		config: cfg,
		db:     db,
	}
}

func (s *Server) Run() error {
	s.setupRoutes()
	return s.router.Run(fmt.Sprintf(":%d", s.config.Port))
}

func (s *Server) setupRoutes() {
	uploadHandler := handlers.NewUploadHandler(s.db)
	s.router.POST("/upload", uploadHandler.HandleUpload)
}
