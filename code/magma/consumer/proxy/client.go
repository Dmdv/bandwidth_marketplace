package proxy

import (
	"github.com/0chain/bandwidth_marketplace/code/pb/consumer"
	"google.golang.org/grpc"
)

// Client creates default Proxy client.
func Client(addr string) (consumer.ProxyClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := consumer.NewProxyClient(conn)

	return client, nil
}
