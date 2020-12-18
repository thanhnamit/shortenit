package internal

import (
	"context"
	"github.com/gorilla/handlers"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/config"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/core"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/platform"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	cfg      *config.Config
	userRepo core.UserRepository
	aliasRepo core.AliasRepository
	aliasSvc core.AliasService
}

func (s *Server) Start() {
	router := NewRouter(s)

	httpSvc := &http.Server{
		Addr: "0.0.0.0:" + s.cfg.Port,
		WriteTimeout: time.Second*15,
		ReadTimeout: time.Second*15,
		IdleTimeout: time.Second*15,
		Handler: handlers.RecoveryHandler()(router.Handler),
	}

	wait := time.Second*15

	go func() {
		if err := httpSvc.ListenAndServe(); err != nil {
			log.Fatalf("Error starting server: %v\n", err)
		}
	}()

	log.Printf("Listening on %s...\n", s.cfg.Port)

	c := make(chan os.Signal, 1)
    // We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
    // SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
    signal.Notify(c, os.Interrupt)

    // Block until we receive our signal.
    <-c

    // Create a deadline to wait for.
    ctx, cancel := context.WithTimeout(context.Background(), wait)
    defer cancel()
    // Doesn't block if no connections, but will otherwise wait
    // until the timeout deadline.
    err := httpSvc.Shutdown(ctx)
    if err != nil {
    	log.Printf("shutting down error :%v", err)
	}
    // Optionally, you could run srv.Shutdown in a goroutine and block on
    // <-ctx.Done() if your application should wait for other services
    // to finalize based on context cancellation.
    log.Println("shutting down")
    os.Exit(0)
}

func NewServer(cfg *config.Config) *Server {
	platform.InitTracer(cfg.TracerName, cfg.TraceCollector)
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