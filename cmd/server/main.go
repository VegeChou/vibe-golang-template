package main

import (
	"log"

	"vibe-golang-template/internal/app"
	"vibe-golang-template/internal/config"
)

func main() {
	cfg := config.Load()
	server := app.NewServer(cfg)

	log.Printf("server starting on %s", cfg.HTTPAddr)
	if err := server.Start(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
