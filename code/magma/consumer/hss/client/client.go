package client

import (
	"github.com/0chain/bandwidth_marketplace/code/pb/consumer"
	"google.golang.org/grpc"
)

// Client creates default HSS client.
func Client(addr string) (consumer.HSSClient, error) {
	conn, err := dial(addr)
	if err != nil {
		return nil, err
	}

	client := consumer.NewHSSClient(conn)

	return client, nil
}

func dial(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, grpc.WithInsecure())
}
