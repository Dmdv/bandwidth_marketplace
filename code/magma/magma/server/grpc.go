package server

import (
	"context"
	"time"

	"github.com/MurashovVen/bandwidth-marketplace/code/core/log"
	"github.com/MurashovVen/bandwidth-marketplace/code/pb/consumer"
	"github.com/MurashovVen/bandwidth-marketplace/code/pb/magma"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	hclient "magma/consumer/hss/client"
	pclient "magma/consumer/proxy/client"
)

func RegisterGRPCServices(server *grpc.Server, preConfiguredConsumers map[string]string) {
	magma.RegisterMagmaServer(server, newMagmaServer(preConfiguredConsumers))
}

type (
	magmaServer struct {
		magma.UnimplementedMagmaServer

		preConfiguredConsumers map[string]string
	}
)

var (
	// Make sure magmaServer implements interface.
	_ magma.MagmaServer = (*magmaServer)(nil)
)

func newMagmaServer(preConfiguredConsumers map[string]string) magma.MagmaServer {
	return &magmaServer{
		preConfiguredConsumers: preConfiguredConsumers,
	}
}

func (s *magmaServer) Connect(ctx context.Context, req *magma.ConnectRequest) (*magma.ConnectResponse, error) {
	log.Logger.Info("Got Connect request")

	consAddr, ok := s.preConfiguredConsumers[req.UserID]
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "consumer with provided ID does not exist")
	}

	err := verifyUser(ctx, req, consAddr)
	if err != nil {
		return nil, err
	}

	err = notifyNewProvider(ctx, req, consAddr)
	if err != nil {
		return nil, err
	}

	// TODO need to implement the next steps according to the provider implementation

	return &magma.ConnectResponse{Status: magma.Status_Connected}, nil
}

func verifyUser(ctx context.Context, req *magma.ConnectRequest, hssAddr string) error {
	cl, err := hclient.Client(hssAddr)
	if err != nil {
		return status.Error(codes.Internal, "can not connect to consumer")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	vuReq := consumer.VerifyUserRequest{
		UserID: req.UserID,
		Auth:   req.Auth,
	}
	resp, err := cl.VerifyUser(ctx, &vuReq)
	if err != nil {
		return status.Errorf(codes.Unknown, "requesting consumer HSS failed with err: %v", err)
	}

	if resp.Status == consumer.VerificationStatus_Unverified {
		return status.Error(codes.Unauthenticated, "user is not authenticated")
	}

	return nil
}

func notifyNewProvider(ctx context.Context, req *magma.ConnectRequest, proxyAddr string) error {
	cl, err := pclient.Client(proxyAddr)
	if err != nil {
		return status.Error(codes.Internal, "can not connect to consumer")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	nnpReq := consumer.NotifyNewProviderRequest{
		SessID:        "", // TODO
		UserID:        req.UserID,
		ProviderID:    req.ProviderID,
		AccessPointID: req.AccessPointID,
	}
	resp, err := cl.NotifyNewProvider(ctx, &nnpReq)
	if err != nil {
		return status.Errorf(codes.Unknown, "requesting consumer Proxy failed with err: %v", err)
	}

	if resp.Status != consumer.ChangeStatus_Success {
		return status.Error(codes.Unknown, "changing provider failed")
	}

	return nil
}
