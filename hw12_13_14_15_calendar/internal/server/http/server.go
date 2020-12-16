package internalhttp

import (
	"context"
	"net"
	"net/http"

	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	httpServer *http.Server
	logger     logger.Logger
}

type ServerConfig struct {
	Host string
	Port string
}

func NewServer(port string, handler http.Handler, l logger.Logger) *Server {
	httpServer := &http.Server{
		Addr:    net.JoinHostPort("", port),
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
