package server

import (
	"context"
	"strconv"

	"github.com/MurashovVen/bandwidth-marketplace/code/core/crypto"
	"github.com/MurashovVen/bandwidth-marketplace/code/core/log"
	"github.com/MurashovVen/bandwidth-marketplace/code/core/slices"
	"github.com/MurashovVen/bandwidth-marketplace/code/core/time"
	"github.com/MurashovVen/bandwidth-marketplace/code/pb/consumer"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"consumer/config"
)

func RegisterGRPCServices(server *grpc.Server, cfg config.HSS) {
	consumer.RegisterHSSServer(server, newHSSServerFromConfig(cfg))
}

type (
	// hssServer represents base implementation of grpc.HSSServer.
	hssServer struct {
		consumer.UnimplementedHSSServer
		users []string // users represents registered users described in config file.
	}
)

var (
	// Make sure hssServer implements interface.
	_ consumer.HSSServer = (*hssServer)(nil)
)

func newHSSServerFromConfig(cfg config.HSS) consumer.HSSServer {
	return &hssServer{
		users: cfg.Users,
	}
}

// VerifyUser checks if the user with the provided id belongs to configured users, verifies signature
// and checks creation date of request.
func (h hssServer) VerifyUser(_ context.Context, request *consumer.VerifyUserRequest) (*consumer.VerifyUserResponse, error) {
	var unverifiedResponse = &consumer.VerifyUserResponse{Status: consumer.VerificationStatus_Unverified}
	log.Logger.Info("Got verifyUser request", zap.Any("request", request))

	switch {
	case !slices.ContainsStr(h.users, request.GetUserID()):
		return unverifiedResponse, status.Error(codes.PermissionDenied, "user with requested id unverified")
	case !isValid(request):
		return unverifiedResponse, status.Error(codes.InvalidArgument, "request message is invalid")
	case !verifySignature(request):
		return unverifiedResponse, status.Error(codes.PermissionDenied, "provided signature is invalid")
	}

	return &consumer.VerifyUserResponse{Status: consumer.VerificationStatus_Verified}, nil
}

// verifySignature verifies signature of creation date from request.
func verifySignature(msg *consumer.VerifyUserRequest) bool {
	hash := crypto.Hash(msg.GetAuth().GetCreationDate())
	ver, err := crypto.Verify(msg.GetUserID(), msg.GetAuth().GetSignature(), hash, msg.GetAuth().GetSignatureScheme())
	return ver && err == nil
}

// isValid checks creation date from request.
func isValid(msg *consumer.VerifyUserRequest) bool {
	ts, err := strconv.Atoi(msg.Auth.CreationDate)
	if err != nil {
		return false
	}

	return int(time.Now()) > ts
}
