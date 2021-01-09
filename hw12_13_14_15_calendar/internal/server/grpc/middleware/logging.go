package middleware

import (
	"context"
	"fmt"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"net"
	"time"
)

type LoggingMiddleware struct {
	logger common.Logger
}

func NewLoggingMiddleware(logger common.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

func (m *LoggingMiddleware) Handle(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	resp, respErr := handler(ctx, req)
	elapsed := time.Since(start)

	p, ok := peer.FromContext(ctx)
	if !ok {
		m.logger.Error("error getting peer")
		return resp, respErr
	}

	ip, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		m.logger.Error("error parsing ip: " + err.Error())
		return resp, respErr
	}

	md, _ := metadata.FromIncomingContext(ctx)

	m.logger.Info(fmt.Sprintf("%s [%s] %s %d %s %s",
		ip,
		start.Format("02/Jan/2006:03:04:05 Z0700"),
		info.FullMethod,
		status.Code(err),
		elapsed,
		md["user-agent"],
	))

	return resp, respErr
}
