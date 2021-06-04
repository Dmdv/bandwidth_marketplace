package context

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/0chain/bandwidth_marketplace/code/core/log"
)

type (
	// CtxKey describe a key used to store values into context.
	CtxKey string
)

var (
	// rootContext represents main context of the app.
	rootContext context.Context

	// rootCancel represents main context cancel function of the app.
	rootCancel context.CancelFunc
)

// SetupRootContext - sets up the root common context and cancel function that can be used to shutdown the node.
func SetupRootContext(ctx context.Context) {
	rootContext, rootCancel = context.WithCancel(ctx)
}

// GetRootContext returns the root common for the server.
// This will be used to control shutting down the server but cleanup all the workers.
func GetRootContext() context.Context {
	return rootContext
}

// Done executes cancel of root context.
func Done() {
	rootCancel()
}

// HandleShutdown handles various shutdown signals.
func HandleShutdown(server *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		for sig := range c {
			switch sig {
			case syscall.SIGINT:
				Done()
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				if err := server.Shutdown(ctx); err != nil {
					log.Logger.Error("Server failed to gracefully shuts down", zap.Error(err))
				}
				cancel()
			case syscall.SIGQUIT:
				Done()
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				if err := server.Shutdown(ctx); err != nil {
					log.Logger.Error("Server failed to gracefully shuts down", zap.Error(err))
				}
				cancel()
			}
		}
	}()
}
