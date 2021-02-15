package internal

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/config"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/core"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/platform"
)

const duration = time.Second * 15

type Server struct {
	cfg       *config.Config
	userRepo  core.UserRepository
	aliasRepo core.AliasRepository
	aliasSvc  core.AliasService
}

func (s *Server) Start() {
	router := NewRouter(s)

	httpSvc := &http.Server{
		Addr:         "0.0.0.0:" + s.cfg.Port,
		WriteTimeout: duration,
		ReadTimeout:  duration,
		IdleTimeout:  duration,
		Handler:      handlers.RecoveryHandler()(router.Handler),
	}

	go func() {
		if err := httpSvc.ListenAndServe(); err != nil {
			log.Fatalf("Error starting server: %v\n", err)
		}
	}()

	log.Printf("Listening on %s...\n", s.cfg.Port)
	log.Printf("CommitSHA %s...\n", s.cfg.CommitSHA)

	s.WaitForInterruptSignal(httpSvc)
}

func (s *Server) WaitForInterruptSignal(httpSvc *http.Server) {
	c := make(chan os.Signal, 1)

	// Accept SIGINT (Ctrl+C)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	err := httpSvc.Shutdown(ctx)
	if err != nil {
		log.Printf("Shutting down error :%v", err)
	}

	log.Println("Shutting down")
	os.Exit(0)
}

func NewServer(cfg *config.Config) *Server {
	platform.InitTracer(cfg.TracerName, cfg.TraceCollector)
	platform.InitMeter()

	ctx := context.Background()
	userRepo := NewUserRepository(ctx, cfg)
	aliasRepo := NewAliasRepository(ctx, cfg)
	aliasSvc := NewAliasClient(cfg)

	return &Server{
		cfg,
		userRepo,
		aliasRepo,
		aliasSvc,
	}
}
