package app

import (
	"fmt"
	"net/http"
	"time"

	"vibe-golang-template/internal/config"
	"vibe-golang-template/internal/controller"
	"vibe-golang-template/internal/i18n"
	"vibe-golang-template/internal/repository/memory"
	"vibe-golang-template/internal/service"
	"vibe-golang-template/pkg/response"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg config.Config) (*Server, error) {
	catalog, err := i18n.LoadCatalog(cfg.I18NFile)
	if err != nil {
		return nil, fmt.Errorf("load i18n config: %w", err)
	}
	response.SetTranslator(catalog)

	userRepo := memory.NewUserRepository()
	userService := service.NewUserService(userRepo)
	api := controller.NewAPI(userService)

	mux := http.NewServeMux()
	api.RegisterRoutes(mux)
	handlerChain := chainMiddlewares(mux, recoverMiddleware)

	return &Server{
		httpServer: &http.Server{
			Addr:              cfg.HTTPAddr,
			Handler:           handlerChain,
			ReadHeaderTimeout: 3 * time.Second,
		},
	}, nil
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
