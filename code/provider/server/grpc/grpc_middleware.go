package grpc

import (
	"context"
	"time"

	"github.com/0chain/bandwidth_marketplace/code/core/datastore"
	"github.com/0chain/bandwidth_marketplace/code/core/errors"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func unaryDatabaseTransactionInjector() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger := ctxzap.Extract(ctx)

		ctx = datastore.GetStore().CreateTransaction(ctx)
		resp, err := handler(ctx, req)
		if err != nil {
			var rollErr = datastore.GetStore().GetTransaction(ctx).Rollback().Error
			if rollErr != nil {
				logger.Error("Couldn't rollback", zap.Error(err))
			}
			return nil, err
		}

		err = datastore.GetStore().GetTransaction(ctx).Commit().Error
		if err != nil {
			code, msg := "commit_error", "error commit to metadata store"
			return nil, errors.WrapError(code, msg, err)
		}

		return resp, err
	}
}

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
			unaryDatabaseTransactionInjector(),
			grpc_recovery.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			unaryTimeoutInterceptor(timeout), // should always be the last, to be "innermost"
		),
	)
}
