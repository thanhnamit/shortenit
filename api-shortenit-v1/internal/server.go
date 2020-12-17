package internal

import (
	"context"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/config"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/core"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/platform"
	"log"
	"net/http"
)

type Server struct {
	cfg      *config.Config
	userRepo core.UserRepository
	aliasSvc core.AliasService
}

func (s *Server) Start() {
	router := NewRouter(s)
	log.Printf("Listening on %s...", s.cfg.Port)
	log.Fatal(http.ListenAndServe(":" + s.cfg.Port, router.Handler))
}

func NewServer(cfg *config.Config) *Server {
	platform.InitTracer(cfg.TracerName, cfg.TraceCollector)
	ctx := context.Background()
	userRepo := NewUserRepository(ctx, cfg)
	aliasSvc := NewAliasClient(cfg)

	return &Server{
		cfg,
		userRepo,
		aliasSvc,
	}
}