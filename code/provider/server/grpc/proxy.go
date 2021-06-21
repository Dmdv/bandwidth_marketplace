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

	"provider/metrics"
)

type (
	proxyServer struct {
		provider.UnimplementedProxyServer

		metrics *metrics.ProxyService
	}
)

func newProxyServer() *proxyServer {
	return &proxyServer{
		metrics: metrics.NewProxyServiceMetrics(),
	}
}

// NewSessionBilling creates Rest Api call to /verifyAcknowledgmentAccepted magma sc endpoint, which verifies Acknowledgment
// Accepting by Consumer with ID from request.
func (p *proxyServer) NewSessionBilling(
	_ context.Context, req *provider.NewSessionBillingRequest) (*provider.NewSessionBillingResponse, error) {
	log.Logger.Debug("Proxy: Got NewSessionBilling request.", zap.Any("request", req))

	params := map[string]string{
		"id":              req.GetAcknowledgmentID(),
		"session_id":      req.GetSessionID(),
		"access_point_id": req.GetAccessPointID(),
		"provider_id":     node.GetSelfNode().ID(),
		"consumer_id":     req.GetConsumerID(),
	}
	_, err := transaction.MakeSCRestAPICall(transaction.MagmaSCAddress, transaction.VerifyAcknowledgmentAcceptedRP, params)
	if err != nil {
		p.metrics.IncAcknowledgmentUnverified()
		return nil, status.Errorf(codes.Unknown, "transaction unverified")
	}

	p.metrics.IncAcknowledgmentVerified()
	log.Logger.Debug("Proxy: Handling NewSessionBilling successfully ended.")
	return &provider.NewSessionBillingResponse{}, nil
}

func (p *proxyServer) ForwardUsage(ctx context.Context, req *provider.ForwardUsageRequest) (*provider.ForwardUsageResponse, error) {
	log.Logger.Debug("Proxy: Got ForwardUsage request.", zap.Any("request", req))

	// TODO need implement after billUsage MagmaSC function implementation

	p.metrics.UpdateDataUploadedMetric(req.GetSessionID(), req.GetOctetsOut())
	p.metrics.UpdateDataDownloadedMetric(req.GetSessionID(), req.GetOctetsIn())
	log.Logger.Debug("Proxy: Handling NewSessionBilling successfully ended.")
	return &provider.ForwardUsageResponse{}, nil
}
