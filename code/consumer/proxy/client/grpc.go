package client

import (
	"github.com/MurashovVen/bandwidth-marketplace/code/pb/consumer"
	"google.golang.org/grpc"
)

// Client creates default grpc.ProxyClient.
func Client(addr string) (consumer.ProxyClient, error) {
	conn, err := dial(addr)
	if err != nil {
		return nil, err
	}

	client := consumer.NewProxyClient(conn)

	return client, nil
}

func dial(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, grpc.WithInsecure())
}
