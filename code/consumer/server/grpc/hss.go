package grpc

import (
	"context"
	"strconv"

	"github.com/0chain/bandwidth_marketplace/code/core/crypto"
	"github.com/0chain/bandwidth_marketplace/code/core/log"
	"github.com/0chain/bandwidth_marketplace/code/core/slices"
	"github.com/0chain/bandwidth_marketplace/code/core/time"
	"github.com/0chain/bandwidth_marketplace/code/pb/consumer"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"consumer/config"
)

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
	log.Logger.Info("HSS: Got VerifyUser request.", zap.Any("request", request))

	var unverifiedResponse = &consumer.VerifyUserResponse{Status: consumer.VerificationStatus_Unverified}
	switch {
	case !slices.ContainsStr(h.users, request.GetUserID()):
		return unverifiedResponse, status.Error(codes.PermissionDenied, "user with requested id unverified")

	case !isValid(request):
		return unverifiedResponse, status.Error(codes.InvalidArgument, "request message is invalid")

	case !verifySignature(request):
		return unverifiedResponse, status.Error(codes.PermissionDenied, "provided signature is invalid")
	}

	log.Logger.Info("HSS: Handling VerifyUser successfully ended.")

	return &consumer.VerifyUserResponse{Status: consumer.VerificationStatus_Verified}, nil
}

// verifySignature verifies signature of creation date from request.
func verifySignature(msg *consumer.VerifyUserRequest) bool {
	var (
		hash      = crypto.Hash(msg.GetAuth().GetCreationDate())
		pbK       = msg.GetAuth().GetPublicKey()
		sign      = msg.GetAuth().GetSignature()
		sigScheme = msg.GetAuth().GetSignatureScheme()
	)
	ver, err := crypto.Verify(pbK, sign, hash, sigScheme)
	return ver && err == nil
}

// isValid checks creation date from request.
func isValid(msg *consumer.VerifyUserRequest) bool {
	ts, err := strconv.Atoi(msg.Auth.CreationDate)

	switch {
	// something went wrong
	case err != nil:
	case int(time.Now()) < ts:
	case msg.GetUserID() != crypto.Hash(msg.GetAuth().GetPublicKey()):

	default: // ok
		return true
	}

	return false
}
