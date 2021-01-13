package internalhttp

import (
	"context"
	"net"
	"net/http"

	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/common"
)

type Server struct {
	httpServer *http.Server
	logger     common.Logger
}

type Config struct {
	Port string
}

func NewServer(cfg Config, handler http.Handler, l common.Logger) *Server {
	httpServer := &http.Server{
		Addr:    net.JoinHostPort("", cfg.Port),
		Handler: loggingMiddleware(handler, l),
	}

	return &Server{
		httpServer: httpServer,
		logger:     l,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
