package grpc

import (
	"context"
	"time"

	"github.com/0chain/bandwidth_marketplace/code/core/log"
	"github.com/0chain/bandwidth_marketplace/code/core/node"
	"github.com/0chain/bandwidth_marketplace/code/pb/magma"
	"github.com/0chain/bandwidth_marketplace/code/pb/provider"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"provider/config"
	mclient "provider/magma"
)

type (
	accessPointServer struct {
		provider.UnimplementedAccessPointServer

		magmaAddress  string
		accessPointID string
	}
)

func newAccessPointServer(cfg *config.Config) *accessPointServer {
	return &accessPointServer{
		magmaAddress:  cfg.MagmaAddress,
		accessPointID: cfg.AccessPointID,
	}
}

// Connect creates Magma Client, makes connect request with provided in provider.ConnectRequest IDs, self node ID,
// and Access Point ID.
//
// Returns success status if execution was successful and no errors.
func (ap *accessPointServer) Connect(ctx context.Context, req *provider.ConnectRequest) (*provider.BeginSession, error) {
	log.Logger.Info("Access Point: Got Connect request.", zap.Any("request", req))

	cl, err := mclient.Client(ap.magmaAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, "can not connect to mclient")
	}

	// connecting to magma
	cReq := magma.ConnectRequest{
		UserID:        req.UserID,
		ConsumerID:    req.ConsumerID,
		ProviderID:    node.GetSelfNode().ID(), // self node represents Provider's node so can use an provider ID
		AccessPointID: ap.accessPointID,
		Auth:          req.Auth,
	}
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	resp, err := cl.Connect(ctx, &cReq)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "requesting Magma with Connect request failed with err: %v", err)
	}

	if resp.Status == magma.ConnectionStatus_Failed {
		return nil, status.Errorf(codes.Unknown, "requesting Magma with Connect request ended with \"Failed\" status")
	}

	log.Logger.Info("Access Point: Handling Connect successfully ended.")

	return &provider.BeginSession{SessionID: resp.SessionID, Status: resp.Status}, nil
}
