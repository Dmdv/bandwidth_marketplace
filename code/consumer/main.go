package main

import (
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/0chain/bandwidth_marketplace/code/core/chain"
	"github.com/0chain/bandwidth_marketplace/code/core/context"
	"github.com/0chain/bandwidth_marketplace/code/core/crypto"
	"github.com/0chain/bandwidth_marketplace/code/core/datastore"
	"github.com/0chain/bandwidth_marketplace/code/core/errors"
	"github.com/0chain/bandwidth_marketplace/code/core/limiter"
	"github.com/0chain/bandwidth_marketplace/code/core/log"
	"github.com/0chain/bandwidth_marketplace/code/core/node"
	"github.com/0chain/bandwidth_marketplace/code/core/transaction"
	"github.com/0chain/gosdk/zcncore"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"consumer/config"
	"consumer/handler"
	"consumer/registration"
	"consumer/server/grpc"
)

func main() {
	cfg := config.Read()
	cfg.CLConfig = config.ExtractCommandLineConfigs()
	log.InitLogging(cfg.Development(), cfg.CLConfig.LogDir, cfg.Logging.Level)

	selfNode := createSelfNode(cfg.CLConfig)

	context.SetupRootContext(node.GetNodeContext())

	datastore.GetStore().OpenWithRetries(cfg.Database.String(), 600)

	serverChain := chain.NewChain(cfg.ServerChain.ID, cfg.ServerChain.OwnerID, cfg.ServerChain.BlockWorker)
	chain.SetServerChain(serverChain)

	registration.Setup(cfg, selfNode.GetWalletString())
	registration.RegisterWithRetries(5)

	pour()

	server := createServer(cfg)
	errMsg := server.ListenAndServe().Error()
	log.Logger.Fatal(errMsg)
}

func pour() {
	log.Logger.Debug("Start pouring wallet")
	txn, err := transaction.NewTransactionEntity()
	if err != nil {
		errors.ExitErr("fail creating faucet transaction", err, 2)
	}

	const (
		address  = "6dba10422e368813802877a85039d3985d96760ed844092319743fb3a76712d3"
		funcName = "pour"
	)

	txnHash, err := txn.ExecuteSmartContract(address, funcName, "", math.MaxInt64)
	if err != nil {
		errors.ExitErr("fail executing pour faucet txn", err, 2)
	}

	_, err = transaction.VerifyTransaction(txnHash)
	if err != nil {
		errors.ExitErr("fail verifying pour faucet txn", err, 2)
	}
	log.Logger.Debug("Success pouring wallet")

	err = zcncore.GetBalance(new(getBalanceCB))
	if err != nil {
		log.Logger.Error("Fail getBalance", zap.Error(err))
	}
}

type getBalanceCB struct{}

func (b *getBalanceCB) OnBalanceAvailable(status int, value int64, info string) {
	log.Logger.Info("Get balance ended.", zap.Int("status", status), zap.Int64("value", value), zap.String("info", info))
}

func createSelfNode(cfg config.CommandLineConfig) *node.Node {
	publicKey, privateKey := crypto.ReadKeysFile(cfg.KeysFile)

	selfNode := node.GetSelfNode()
	selfNode.SetKeys(publicKey, privateKey)
	selfNode.SetHostURL(cfg.Host, cfg.Port)
	log.Logger.Info("Base URL: " + selfNode.GetURLBase())

	if err := selfNode.Validate(); err != nil {
		log.Logger.Panic("invalid self node", zap.Error(err))
	}

	log.Logger.Info("Self identity", zap.String("id", selfNode.ID()))

	return selfNode
}

func createServer(cfg *config.Config) (server *http.Server) {
	// setup CORS
	router := mux.NewRouter()
	limiter.ConfigRateLimits(cfg.Handler.RateLimit)
	handler.SetupHandlers(router)

	address := ":" + strconv.Itoa(cfg.CLConfig.Port)
	originsOk := handlers.AllowedOriginValidator(isValidOrigin)
	headersOk := handlers.AllowedHeaders([]string{
		"X-Requested-With", "X-App-Client-ID",
		"X-App-Client-Key", "Content-Type",
	})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	server = &http.Server{
		Addr:              address,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		MaxHeaderBytes:    1 << 20,
		Handler:           handlers.CORS(originsOk, headersOk, methodsOk)(router),
	}
	if cfg.Development() { // non idle & write timeouts setup to enable pprof
		server.IdleTimeout = 0
		server.WriteTimeout = 0
	}

	context.HandleShutdown(server)
	handler.HandleShutdown(context.GetRootContext())
	node.Start()
	log.Logger.Info("Ready to listen to the requests")

	go startGRPC(cfg)

	return server
}

func startGRPC(cfg *config.Config) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.CLConfig.GrpcPort))
	if err != nil {
		log.Logger.Fatal("failed to listen", zap.Error(err))
	}

	timeout := time.Duration(cfg.GRPCServerTimeout) * time.Second
	server := grpc.NewServerWithMiddlewares(log.Logger, timeout)
	grpc.RegisterGRPCServices(server, cfg)

	errMsg := server.Serve(listener).Error()
	log.Logger.Fatal(errMsg)
}

func isValidOrigin(origin string) bool {
	uri, err := url.Parse(origin)
	if err != nil {
		return false
	}

	host := uri.Hostname()
	switch { // allowed origins
	case host == "localhost":
	case host == "0chain.net":
	case strings.HasSuffix(host, ".0chain.net"):
	case strings.HasSuffix(host, ".alphanet-0chain.net"):
	case strings.HasSuffix(host, ".devnet-0chain.net"):
	case strings.HasSuffix(host, ".testnet-0chain.net"):
	case strings.HasSuffix(host, ".mainnet-0chain.net"):

	default: // not allowed
		return false
	}

	return true
}
