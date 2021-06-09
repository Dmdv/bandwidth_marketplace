package grpc

import (
	"github.com/0chain/bandwidth_marketplace/code/pb/provider"
	"google.golang.org/grpc"

	"provider/config"
)

func RegisterGRPCServices(server *grpc.Server, cfg *config.Config) {
	provider.RegisterProxyServer(server, newProxyServer())
	provider.RegisterAccessPointServer(server, newAccessPointServer(cfg))
}
