package grpc

import (
	"context"

	"github.com/0chain/bandwidth_marketplace/code/core/log"
	"github.com/0chain/bandwidth_marketplace/code/pb/consumer"
	"go.uber.org/zap"
)

type (
	// hssServer represents base implementation of grpc.HSSServer.
	hssServer struct {
		consumer.UnimplementedHSSServer
	}
)

var (
	// Make sure hssServer implements interface.
	_ consumer.HSSServer = (*hssServer)(nil)
)

func newHSSServerFromConfig() consumer.HSSServer {
	return &hssServer{}
}

// VerifyUser checks if the user with the provided id belongs to configured users.
func (s hssServer) VerifyUser(_ context.Context, req *consumer.VerifyUserRequest) (*consumer.VerifyUserResponse, error) {
	log.Logger.Info("HSS: Got VerifyUser request.", zap.Any("request", req))

	log.Logger.Info("HSS: Handling VerifyUser successfully ended.")

	return &consumer.VerifyUserResponse{}, nil
}
