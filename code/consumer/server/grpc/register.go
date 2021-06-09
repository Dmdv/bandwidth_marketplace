package grpc

import (
	"github.com/0chain/bandwidth_marketplace/code/pb/consumer"
	"google.golang.org/grpc"

	"consumer/config"
)

func RegisterGRPCServices(server *grpc.Server, cfg *config.Config) {
	consumer.RegisterHSSServer(server, newHSSServerFromConfig(cfg.HSS))
	consumer.RegisterProxyServer(server, newProxyServer(cfg.Proxy, cfg.MagmaAddress))
}
