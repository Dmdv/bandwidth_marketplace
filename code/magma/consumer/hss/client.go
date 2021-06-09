package hss

import (
	"github.com/0chain/bandwidth_marketplace/code/pb/consumer"
	"google.golang.org/grpc"
)

// Client creates default HSS client.
func Client(addr string) (consumer.HSSClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := consumer.NewHSSClient(conn)

	return client, nil
}
