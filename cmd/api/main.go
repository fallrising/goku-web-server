package main

import (
	"log"

	v1 "github.com/fallrising/goku-api/api/v1"
	"github.com/fallrising/goku-api/internal/database"
	"github.com/fallrising/goku-api/internal/server"
	"github.com/fallrising/goku-api/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.NewDatabase(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	v1Router := v1.NewRouter(db)

	srv := server.NewServer(cfg, v1Router)
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
