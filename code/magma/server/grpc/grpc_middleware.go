package grpc

import (
	"context"
	"time"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func unaryTimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return handler(ctx, req)
	}
}

func NewServerWithMiddlewares(logger *zap.Logger, timeout time.Duration) *grpc.Server {
	return grpc.NewServer(
		grpc.ChainStreamInterceptor(
			grpc_zap.StreamServerInterceptor(logger),
			grpc_recovery.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
		),
		grpc.ChainUnaryInterceptor(
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_recovery.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			unaryTimeoutInterceptor(timeout), // should always be the last, to be "innermost"
		),
	)
}
