package internalgrpc

import (
	"net"

	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/api"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/server/grpc/middleware"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer   *grpc.Server
	port         string
	eventService api.EventServiceServer
	middlewares  []middleware.Middleware
}

func NewServer(
	port string,
	eventService api.EventServiceServer,
	middlewares []middleware.Middleware,
) (*Server, error) {
	return &Server{
		port:         port,
		eventService: eventService,
		middlewares:  middlewares,
	}, nil
}

func (s *Server) Start() error {
	addr := net.JoinHostPort("", s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	interceptors := make([]grpc.UnaryServerInterceptor, 0, len(s.middlewares))
	for _, mdlw := range s.middlewares {
		interceptors = append(interceptors, mdlw.Handle)
	}
	s.grpcServer = grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors...))

	api.RegisterEventServiceServer(s.grpcServer, s.eventService)

	if err := s.grpcServer.Serve(listener); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}
