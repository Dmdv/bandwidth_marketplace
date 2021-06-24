package grpc

import (
	"context"
	"encoding/json"

	"github.com/0chain/bandwidth_marketplace/code/core/log"
	"github.com/0chain/bandwidth_marketplace/code/core/provider"
	"github.com/0chain/bandwidth_marketplace/code/core/transaction"
	"github.com/0chain/bandwidth_marketplace/code/pb/consumer"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"consumer/config"
)

type (
	proxyServer struct {
		consumer.UnimplementedProxyServer

		cfg          *config.Proxy
		magmaAddress string
	}
)

var (
	// Make sure proxyServer implements interface.
	_ consumer.ProxyServer = (*proxyServer)(nil)
)

func newProxyServer(cfg config.Proxy, magmaAddress string) consumer.ProxyServer {
	return &proxyServer{
		cfg:          &cfg,
		magmaAddress: magmaAddress,
	}
}

// NotifyNewProvider starts picking provider process and returns consumer.ChangeProviderResponse.
func (p proxyServer) NotifyNewProvider(_ context.Context, request *consumer.NotifyNewProviderRequest) (*consumer.ChangeProviderResponse, error) {
	log.Logger.Info("Proxy: Got NotifyNewProvider request", zap.Any("request", request))

	acknID, err := pickProvider(request, p.cfg.Terms)
	if err != nil {
		return nil, err
	}

	log.Logger.Info("Proxy: Handling NotifyNewProvider successfully ended.")

	return &consumer.ChangeProviderResponse{Status: consumer.ChangeStatus_Success, AcknowledgmentID: acknID}, nil
}

// pickProvider requests to Magma Smart Contract for provider.Terms of Provider with ID passed in args,
// checks gotten terms for validity with config.Terms and decides changing provider with provider.ExecuteAcceptTerms or not.
//
// Returns Acknowledgment ID.
//
// If an error occurs while executing, Magma service with provided address will be notified with grpc.Status_Failed.
// If executing will be ended successfully, Magma will be notified with grpc.Status_Success.
func pickProvider(req *consumer.NotifyNewProviderRequest, cfg config.Terms) (string, error) {
	log.Logger.Info("Got notifyNewProvider request, trying to pick provider ...")

	// extract terms from block workers by executing sc api
	terms, err := respondTerms(req.ProviderID)
	if err != nil {
		log.Logger.Error("Responding terms failed", zap.Error(err))

		return "", status.Errorf(codes.Unknown, "responding terms of provider failed with err: %v", err)
	}

	if err := cfg.Validate(terms); err != nil {
		log.Logger.Error("Terms is invalid", zap.Any("terms", terms), zap.Error(err))

		return "", status.Error(codes.InvalidArgument, "terms can not be picked cause of configured requirements")
	}

	// accept terms
	ackn, err := provider.ExecuteAcceptTerms(req.ProviderID, req.AccessPointID, req.SessID)
	if err != nil {
		log.Logger.Error("Accepting terms failed", zap.Error(err))

		return "", status.Errorf(codes.Unknown, "accepting terms of provider failed with err: %v", err)
	}

	log.Logger.Info("Picking provider ended successfully", zap.Any("acknowledgment", ackn))

	return ackn.SessionID, nil
}

// respondTerms responds provider.Terms from blockchain.
func respondTerms(providerID string) (*provider.ProviderTerms, error) {
	params := map[string]string{
		"provider_id": providerID,
	}
	res, err := transaction.MakeSCRestAPICall(transaction.MagmaSCAddress, transaction.ProviderTermsRP, params)
	if err != nil {
		return nil, err
	}

	terms := new(provider.ProviderTerms)
	err = json.Unmarshal(res, terms)
	if err != nil {
		return nil, err
	}

	return terms, err
}
