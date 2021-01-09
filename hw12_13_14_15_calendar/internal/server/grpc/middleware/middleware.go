package middleware

import (
	"context"
	"google.golang.org/grpc"
)

type Middleware interface {
	Handle(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error)
}
