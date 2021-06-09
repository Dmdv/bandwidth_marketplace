package magma

import (
	"github.com/0chain/bandwidth_marketplace/code/pb/magma"
	"google.golang.org/grpc"
)

// Client creates default Magma client.
func Client(addr string) (magma.MagmaClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := magma.NewMagmaClient(conn)

	return client, nil
}
