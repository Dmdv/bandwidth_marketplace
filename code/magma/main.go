package main

import (
	ctx "context"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/MurashovVen/bandwidth-marketplace/code/core/context"
	"github.com/MurashovVen/bandwidth-marketplace/code/core/limiter"
	"github.com/MurashovVen/bandwidth-marketplace/code/core/log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"magma/config"
	"magma/handler"
	mserver "magma/magma/server"
)

func main() {
	cfg := config.Read()

	log.InitLogging(true, cfg.Logging.Dir, cfg.Logging.Level)

	context.SetupRootContext(ctx.Background())

	server := createServer(cfg)
	errMsg := server.ListenAndServe().Error()
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
	log.Logger.Info("Ready to listen to the requests")

	go startGRPC(cfg)

	return server
}

func startGRPC(cfg *config.Config) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.GRPCPort))
	if err != nil {
		log.Logger.Fatal("failed to listen", zap.Error(err))
	}

	server := handler.NewServerWithMiddlewares(log.Logger)

	// key - ID, value - address
	preConfiguredConsumers := map[string]string{
		"3d8b2c47542404b0059a4c320dca873645d0a3259532a8e17e4cd10db2ef2012b0ba6ef46c5dffab0bb43c05cfc3da7004adf5fe5497e339b298f742fa9843a3": ":7031",
	}
	mserver.RegisterGRPCServices(server, preConfiguredConsumers)

	log.Logger.Info("Server preparing to be started")
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
