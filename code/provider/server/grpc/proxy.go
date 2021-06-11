package grpc

import (
	"context"

	"github.com/0chain/bandwidth_marketplace/code/core/log"
	"github.com/0chain/bandwidth_marketplace/code/core/node"
	"github.com/0chain/bandwidth_marketplace/code/core/transaction"
	"github.com/0chain/bandwidth_marketplace/code/pb/provider"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	proxyServer struct {
		provider.UnimplementedProxyServer
	}
)

func newProxyServer() *proxyServer {
	return &proxyServer{}
}

// NewSessionBilling creates Rest Api call to /verifyAcknowledgmentAccepted magma sc endpoint, which verifies Acknowledgment
// Accepting by Consumer with ID from request.
func (p *proxyServer) NewSessionBilling(ctx context.Context, req *provider.NewSessionBillingRequest) (*provider.NewSessionBillingResponse, error) {
	log.Logger.Info("Proxy: Got NewSessionBilling request.", zap.Any("request", req))

	params := map[string]string{
		"id":              req.AcknowledgmentID,
		"session_id":      req.SessionID,
		"access_point_id": req.AccessPointID,
		"provider_id":     node.GetSelfNode().ID(),
		"consumer_id":     req.ConsumerID,
	}
	_, err := transaction.MakeSCRestAPICall(transaction.MagmaSCAddress, transaction.VerifyAcknowledgmentAcceptedRP, params)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "transaction unverified")
	}

	log.Logger.Info("Proxy: Handling NewSessionBilling successfully ended.")

	return &provider.NewSessionBillingResponse{}, nil
}

func (p *proxyServer) ForwardUsage(ctx context.Context, req *provider.ForwardUsageRequest) (*provider.ForwardUsageResponse, error) {
	log.Logger.Info("Proxy: Got ForwardUsage request.", zap.Any("request", req))

	// TODO need implement after billUsage MagmaSC function implementation

	log.Logger.Info("Proxy: Handling NewSessionBilling successfully ended.")

	return &provider.ForwardUsageResponse{}, nil
}
