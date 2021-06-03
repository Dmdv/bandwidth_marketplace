package handler

import (
	"context"
	"time"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/MurashovVen/bandwidth-marketplace/code/core/datastore"
	"github.com/MurashovVen/bandwidth-marketplace/code/core/errors"
)

const (
	// TimeoutSeconds represents the deadline for handling requests.
	TimeoutSeconds = 10
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
			unaryDatabaseTransactionInjector(),
			grpc_recovery.UnaryServerInterceptor(),
			unaryTimeoutInterceptor(), // should always be the last, to be "innermost"
		),
	)
}
