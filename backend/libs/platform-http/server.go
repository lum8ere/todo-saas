package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	platformlogger "github.com/lum1ere/todo-saas/backend/libs/platform-logger"
	"go.uber.org/zap"
)

type ServerConfig struct {
	Addr            string        // ":8080"
	ShutdownTimeout time.Duration // 10 * time.Second
}

// Server — враппер над http.Server с логгером и graceful shutdown.
type Server struct {
	httpServer *http.Server
	logger     *platformlogger.Logger
}

func NewServer(cfg ServerConfig, router chi.Router, logger *platformlogger.Logger) *Server {
	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = 10 * time.Second
	}

	s := &http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	return &Server{
		httpServer: s,
		logger:     logger,
	}
}

// Run блокирует goroutine, пока не придёт SIGINT/SIGTERM.
func (s *Server) Run() {
	log := s.logger.Base
	log.Info("starting http server", zap.String("addr", s.httpServer.Addr))

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("http server failed", zap.Error(err))
		}
	}()

	// Ждём сигнала остановки
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("shutting down http server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Error("graceful shutdown failed", zap.Error(err))
	} else {
		log.Info("http server stopped")
	}
}
