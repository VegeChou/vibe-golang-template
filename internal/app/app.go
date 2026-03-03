package app

import (
	"net/http"
	"time"

	"vibe-golang-template/internal/config"
	"vibe-golang-template/internal/handler"
	"vibe-golang-template/internal/repository/memory"
	"vibe-golang-template/internal/service"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg config.Config) *Server {
	userRepo := memory.NewUserRepository()
	userService := service.NewUserService(userRepo)
	api := handler.NewAPI(userService)

	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	return &Server{
		httpServer: &http.Server{
			Addr:              cfg.HTTPAddr,
			Handler:           mux,
			ReadHeaderTimeout: 3 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
