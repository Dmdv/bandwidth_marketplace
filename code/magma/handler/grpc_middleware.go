package handler

import (
	"context"
	"time"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	// TimeoutSeconds represents the deadline for handling requests.
	TimeoutSeconds = 10
)

func unaryTimeoutInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		deadline := time.Now().Add(TimeoutSeconds * time.Second)
		ctx, cancel := context.WithDeadline(ctx, deadline)
		defer cancel()
		return handler(ctx, req)
	}
}

func NewServerWithMiddlewares(logger *zap.Logger) *grpc.Server {
	return grpc.NewServer(
		grpc.ChainStreamInterceptor(
			grpc_zap.StreamServerInterceptor(logger),
			grpc_recovery.StreamServerInterceptor(),
		),
		grpc.ChainUnaryInterceptor(
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_recovery.UnaryServerInterceptor(),
			unaryTimeoutInterceptor(), // should always be the last, to be "innermost"
		),
	)
}
