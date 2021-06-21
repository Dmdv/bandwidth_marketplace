package grpc

import (
	"context"
	"time"

	"github.com/0chain/bandwidth_marketplace/code/core/log"
	"github.com/0chain/bandwidth_marketplace/code/pb/consumer"
	"github.com/0chain/bandwidth_marketplace/code/pb/provider"
	magma "github.com/magma/augmented-networks/accounting/protos"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/config"
	"magma/consumer/hss"
	cproxy "magma/consumer/proxy"
	"magma/metrics"
	pproxy "magma/provider/proxy"
)

func RegisterGRPCServices(server *grpc.Server, cfg *config.Config) {
	magma.RegisterAccountingServer(server, newMagmaServer(cfg))
}

type (
	magmaServer struct {
		magma.UnimplementedAccountingServer

		consumerAddress      string
		providerAddress      string
		defaultClientTimeout time.Duration
		idMaps               config.IDMaps

		metrics *metrics.MagmaService
	}
)

var (
	// Make sure magmaServer implements interface.
	_ magma.AccountingServer = (*magmaServer)(nil)
)

func newMagmaServer(cfg *config.Config) magma.AccountingServer {
	return &magmaServer{
		consumerAddress:      cfg.ConsumerAddress,
		providerAddress:      cfg.ProviderAddress,
		defaultClientTimeout: time.Duration(cfg.GRPCDefaultClientTimeout) * time.Second,
		idMaps:               cfg.IDMaps,
		metrics:              metrics.NewMagmaServiceMetrics(),
	}
}

// Start makes verification request to Consumer.HSS service,
// notifies Consumer.Proxy with NotifyNewProvider request
// and notifies Provider.Proxy with NewSessionBilling request.
func (s *magmaServer) Start(ctx context.Context, req *magma.Session) (*magma.SessionResp, error) {
	log.Logger.Debug("Magma: Got Start request", zap.Any("request", req))

	if err := verifyUser(ctx, req, s.consumerAddress); err != nil {
		return nil, err
	}
	log.Logger.Debug("Magma.Start: Consumer.HSS.VerifyUser successfully executed.")

	acknID, err := s.notifyNewProvider(ctx, req)
	if err != nil {
		return nil, err
	}
	log.Logger.Debug("Magma.Start: Consumer.Proxy.NotifyNewProvider successfully executed.")

	if err := s.newSessionBilling(ctx, req, acknID); err != nil {
		return nil, err
	}
	log.Logger.Debug("Magma.Start: Provider.Proxy.NewSessionBilling successfully executed.")

	s.metrics.IncSessionStarted()
	log.Logger.Debug("Magma: Handling Start successfully ended.")
	return &magma.SessionResp{}, nil
}

// verifyUser creates Consumer HSS client and makes VerifyUser request.
//
// It returns already configured status error with codes if error occurs while execution.
func verifyUser(ctx context.Context, req *magma.Session, hssAddr string) error {
	cl, err := hss.Client(hssAddr)
	if err != nil {
		return status.Error(codes.Internal, "can not connect to consumer HSS making VerifyUser request")
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	vuReq := consumer.VerifyUserRequest{
		UserID: req.GetName(),
	}
	_, err = cl.VerifyUser(ctx, &vuReq)
	if err != nil {
		f := "requesting consumer HSS with VerifyUser request failed with err: %v"
		return status.Errorf(codes.Unknown, f, err)
	}

	return nil
}

// notifyNewProvider creates Consumer Proxy client and makes NotifyNewProvider request with IDs from provided
// session request and provided session ID.
//
// It returns already configured status error with codes if error occurs while execution,
// else returns Acknowledgment ID.
func (s *magmaServer) notifyNewProvider(ctx context.Context, req *magma.Session) (string, error) {
	cl, err := cproxy.Client(s.consumerAddress)
	if err != nil {
		msg := "can not connect to consumer Proxy while making NotifyNewProvider request"
		return "", status.Error(codes.Internal, msg)
	}

	ctx, cancel := context.WithTimeout(ctx, s.defaultClientTimeout)
	defer cancel()
	nnpReq := consumer.NotifyNewProviderRequest{
		SessID:        req.GetSessionId(),
		UserID:        req.GetIMSI(),
		ProviderID:    s.idMaps.Provider[req.GetProviderId()],
		AccessPointID: req.GetProviderApn(),
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
func (s *magmaServer) newSessionBilling(ctx context.Context, req *magma.Session, acknowledgmentID string) error {
	cl, err := pproxy.Client(s.providerAddress)
	if err != nil {
		msg := "can not connect to provider Proxy while making NewSessionBilling request"
		return status.Error(codes.Internal, msg)
	}

	ctx, cancel := context.WithTimeout(ctx, s.defaultClientTimeout)
	defer cancel()
	nsbReq := provider.NewSessionBillingRequest{
		SessionID:        req.GetSessionId(),
		UserID:           req.GetIMSI(),
		ConsumerID:       s.idMaps.Consumer[req.GetConsumerId()],
		AccessPointID:    req.GetProviderApn(),
		AcknowledgmentID: acknowledgmentID,
	}
	_, err = cl.NewSessionBilling(ctx, &nsbReq)
	if err != nil {
		f := "requesting provider Proxy with NewSessionBilling request failed with err: %v"
		return status.Errorf(codes.Unknown, f, err)
	}

	return nil
}

// Update forwards update request from Magma to Provider.Proxy with ForwardUsage request.
func (s *magmaServer) Update(ctx context.Context, req *magma.UpdateReq) (*magma.SessionResp, error) {
	log.Logger.Debug("Magma: Got Update request.", zap.Any("request", req))

	cl, err := pproxy.Client(s.providerAddress)
	if err != nil {
		msg := "can not connect to provider Proxy while making ForwardUsage request"
		return nil, status.Error(codes.Internal, msg)
	}

	fuReq := provider.ForwardUsageRequest{
		SessionID:   req.GetSession().GetSessionId(),
		OctetsIn:    req.GetOctetsIn(),
		OctetsOut:   req.GetOctetsOut(),
		SessionTime: req.GetSessionTime(),
	}
	ctx, cancel := context.WithTimeout(ctx, s.defaultClientTimeout)
	defer cancel()
	if _, err = cl.ForwardUsage(ctx, &fuReq); err != nil {
		f := "requesting consumer Proxy with NotifyNewProvider request failed with err: %v"
		return nil, status.Errorf(codes.Unknown, f, err)
	}

	s.metrics.UpdateDataUploadedMetric(req.GetSession().GetSessionId(), req.GetOctetsOut())
	s.metrics.UpdateDataDownloadedMetric(req.GetSession().GetSessionId(), req.GetOctetsIn())
	log.Logger.Debug("Magma: Handling Update successfully ended.")
	return &magma.SessionResp{}, err
}

// Stop executes Magma.Update procedure.
// NOTE: need to be modified.
func (s magmaServer) Stop(ctx context.Context, req *magma.UpdateReq) (*magma.StopResp, error) {
	log.Logger.Debug("Magma: Got Stop request.", zap.Any("request", req))

	_, err := s.Update(ctx, req)
	if err != nil {
		return nil, err
	}
	// TODO implement magma sc function

	s.metrics.IncSessionStopped()
	log.Logger.Debug("Magma: Handling Stop successfully ended.")
	return &magma.StopResp{}, nil
}
