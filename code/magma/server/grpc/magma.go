package grpc

import (
	"context"
	"strconv"
	"time"

	"github.com/0chain/bandwidth_marketplace/code/core/crypto"
	"github.com/0chain/bandwidth_marketplace/code/core/log"
	"github.com/0chain/bandwidth_marketplace/code/pb/consumer"
	"github.com/0chain/bandwidth_marketplace/code/pb/magma"
	"github.com/0chain/bandwidth_marketplace/code/pb/provider"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/consumer/hss"
	cproxy "magma/consumer/proxy"
	pproxy "magma/provider/proxy"
)

func RegisterGRPCServices(server *grpc.Server, preConfiguredConsumers, preConfiguredProviders map[string]string) {
	magma.RegisterMagmaServer(server, newMagmaServer(preConfiguredConsumers, preConfiguredProviders))
}

type (
	magmaServer struct {
		magma.UnimplementedMagmaServer

		preConfiguredConsumers map[string]string
		preConfiguredProviders map[string]string
	}
)

var (
	// Make sure magmaServer implements interface.
	_ magma.MagmaServer = (*magmaServer)(nil)
)

func newMagmaServer(preConfiguredConsumers, preConfiguredProviders map[string]string) magma.MagmaServer {
	return &magmaServer{
		preConfiguredConsumers: preConfiguredConsumers,
		preConfiguredProviders: preConfiguredProviders,
	}
}

// Connect handles connection request to Magma. It verifies user by making VerifyUser request to pre configured Consumer HSS,
// notifies pre configured Provider Proxy by making NotifyNewProvider request and notifies Provider Proxy with
// NewSessionBilling request.
func (s *magmaServer) Connect(ctx context.Context, req *magma.ConnectRequest) (*magma.NewSessionStart, error) {
	log.Logger.Info("Magma: Got Connect request", zap.Any("request", req))

	consAddr, ok := s.preConfiguredConsumers[req.UserID]
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "consumer with provided ID does not exist")
	}

	provAddr, ok := s.preConfiguredProviders[req.ProviderID]
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "provider with provided ID does not exist")
	}

	if err := verifyUser(ctx, req, consAddr); err != nil {
		return nil, err
	}

	sessionID := crypto.Hash(strconv.Itoa(int(time.Now().UnixNano()))) // TODO change on something better
	acknID, err := notifyNewProvider(ctx, req, consAddr, sessionID)
	if err != nil {
		return nil, err
	}

	if err := newSessionBilling(ctx, req, acknID, sessionID, provAddr); err != nil {
		return nil, err
	}

	log.Logger.Info("Magma: Handling Connect successfully ended.")

	return &magma.NewSessionStart{Status: magma.ConnectionStatus_Connected, SessionID: sessionID}, nil
}

// verifyUser creates Consumer HSS client and makes VerifyUser request with user ID and auth data from provided
// magma.ConnectRequest.
//
// It returns already configured status error with codes if error occurs while execution.
func verifyUser(ctx context.Context, req *magma.ConnectRequest, hssAddr string) error {
	cl, err := hss.Client(hssAddr)
	if err != nil {
		return status.Error(codes.Internal, "can not connect to consumer HSS making VerifyUser request")
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	vuReq := consumer.VerifyUserRequest{
		UserID: req.UserID,
		Auth:   req.Auth,
	}
	resp, err := cl.VerifyUser(ctx, &vuReq)
	if err != nil {
		f := "requesting consumer HSS with VerifyUser request failed with err: %v"
		return status.Errorf(codes.Unknown, f, err)
	}

	if resp.Status == consumer.VerificationStatus_Unverified {
		return status.Error(codes.Unauthenticated, "user is not authenticated")
	}

	return nil
}

// notifyNewProvider creates Consumer Proxy client and makes NotifyNewProvider request with IDs from provided
// magma.ConnectRequest and provided session ID.
//
// It returns already configured status error with codes if error occurs while execution,
// else returns Acknowledgment ID.
func notifyNewProvider(ctx context.Context, req *magma.ConnectRequest, proxyAddr, sessionID string) (string, error) {
	cl, err := cproxy.Client(proxyAddr)
	if err != nil {
		msg := "can not connect to consumer Proxy while making NotifyNewProvider request"
		return "", status.Error(codes.Internal, msg)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	nnpReq := consumer.NotifyNewProviderRequest{
		SessID:        sessionID,
		UserID:        req.UserID,
		ProviderID:    req.ProviderID,
		AccessPointID: req.AccessPointID,
	}
	resp, err := cl.NotifyNewProvider(ctx, &nnpReq)
	if err != nil {
		f := "requesting consumer Proxy with NotifyNewProvider request failed with err: %v"
		return "", status.Errorf(codes.Unknown, f, err)
	}

	if resp.Status != consumer.ChangeStatus_Success {
		return "", status.Error(codes.Unknown, "changing provider failed")
	}

	return resp.AcknowledgmentID, nil
}

// newSessionBilling creates Provider Proxy client and makes NewSessionBilling request with IDs from provided
// magma.ConnectRequest, session ID and Acknowledgment ID.
//
// It returns already configured status error with codes if error occurs while execution.
func newSessionBilling(ctx context.Context, req *magma.ConnectRequest, acknowledgmentID, sessionID, proxyAddr string) error {
	cl, err := pproxy.Client(proxyAddr)
	if err != nil {
		msg := "can not connect to provider Proxy while making NewSessionBilling request"
		return status.Error(codes.Internal, msg)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	nsbReq := provider.NewSessionBillingRequest{
		SessionID:        sessionID,
		UserID:           req.UserID,
		ConsumerID:       req.ConsumerID,
		AccessPointID:    req.AccessPointID,
		AcknowledgmentID: acknowledgmentID,
	}
	_, err = cl.NewSessionBilling(ctx, &nsbReq)
	if err != nil {
		f := "requesting provider Proxy with NewSessionBilling request failed with err: %v"
		return status.Errorf(codes.Unknown, f, err)
	}

	return nil
}

// ReportUsage handles reporting request from Provider's Access Point
// and forwards it to pre configured Provider's Proxy.
func (s *magmaServer) ReportUsage(ctx context.Context, req *magma.ReportUsageRequest) (*magma.ReportUsageResponse, error) {
	log.Logger.Info("Magma: Got ReportUsage request.", zap.Any("request", req))

	provAddr, ok := s.preConfiguredProviders[req.ProviderID]
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "provider with provided ID does not exist")
	}

	cl, err := pproxy.Client(provAddr)
	if err != nil {
		msg := "can not connect to provider Proxy while making ForwardUsage request"
		return nil, status.Error(codes.Internal, msg)
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if _, err = cl.ForwardUsage(ctx, req.UsageData); err != nil {
		f := "requesting consumer Proxy with NotifyNewProvider request failed with err: %v"
		return nil, status.Errorf(codes.Unknown, f, err)
	}

	log.Logger.Info("Magma: Handling ReportUsage successfully ended.")

	return &magma.ReportUsageResponse{}, err
}
