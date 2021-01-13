package middleware

import (
	"context"

	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/user"
	handler2 "github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/server/grpc/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthenticationMiddleware struct {
}

func NewAuthenticationMiddleware() *AuthenticationMiddleware {
	return &AuthenticationMiddleware{}
}

func (m *AuthenticationMiddleware) Handle(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if len(md["x-uid"]) == 0 {
		return nil, status.Error(codes.Unauthenticated, "pass x-uid parameter in metadata")
	}

	uid := user.UID(md["x-uid"][0])
	authCtx := context.WithValue(ctx, handler2.UIDKey{}, uid)

	return handler(authCtx, req)
}
