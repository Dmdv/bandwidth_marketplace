package main

import (
	ctx "context"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/0chain/bandwidth_marketplace/code/core/context"
	"github.com/0chain/bandwidth_marketplace/code/core/limiter"
	"github.com/0chain/bandwidth_marketplace/code/core/log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"magma/config"
	"magma/handler"
	"magma/server/grpc"
)

func main() {
	cfg := config.Read()

	log.InitLogging(true, cfg.Logging.Dir, cfg.Logging.Level)

	context.SetupRootContext(ctx.Background())

	serv := createServer(cfg)
	errMsg := serv.ListenAndServe().Error()
	log.Logger.Fatal(errMsg)
}

func createServer(cfg *config.Config) (server *http.Server) {
	// setup CORS
	router := mux.NewRouter()
	limiter.ConfigRateLimits(cfg.Handler.RateLimit)
	handler.SetupHandlers(router)

	address := ":" + strconv.Itoa(cfg.Port)
	originsOk := handlers.AllowedOriginValidator(isValidOrigin)
	headersOk := handlers.AllowedHeaders([]string{
		"X-Requested-With", "X-App-Client-ID",
		"X-App-Client-Key", "Content-Type",
	})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	server = &http.Server{
		Addr:              address,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    1 << 20,
		Handler:           handlers.CORS(originsOk, headersOk, methodsOk)(router),
	}

	context.HandleShutdown(server)

	go startGRPC(cfg)

	return server
}

func startGRPC(cfg *config.Config) {
	listener, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		log.Logger.Fatal("failed to listen", zap.Error(err))
	}

	timeout := time.Duration(cfg.GRPCServerTimeout) * time.Second
	serv := grpc.NewServerWithMiddlewares(log.Logger, timeout)
	grpc.RegisterGRPCServices(serv, cfg)

	log.Logger.Info("GRPC server is preparing to be started.", zap.String("grps address", cfg.GRPCAddress))
	errMsg := serv.Serve(listener).Error()
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
