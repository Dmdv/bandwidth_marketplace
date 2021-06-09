package proxy

import (
	"github.com/0chain/bandwidth_marketplace/code/pb/provider"
	"google.golang.org/grpc"
)

// Client creates default Proxy client.
func Client(addr string) (provider.ProxyClient, error) {
	conn, err := dial(addr)
	if err != nil {
		return nil, err
	}

	client := provider.NewProxyClient(conn)

	return client, nil
}

func dial(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, grpc.WithInsecure())
}
