package grpc

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/0chain/bandwidth_marketplace/code/core/crypto"
	"github.com/0chain/bandwidth_marketplace/code/pb/consumer"
	"github.com/0chain/bandwidth_marketplace/code/pb/provider"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const (
	pbK  = "3d8b2c47542404b0059a4c320dca873645d0a3259532a8e17e4cd10db2ef2012b0ba6ef46c5dffab0bb43c05cfc3da7004adf5fe5497e339b298f742fa9843a3"
	sign = "a2338d9493bed993bab028789cd61f47594614567218db1bbe70b0900d578097"
)

func Test_accessPointServer_Connect(t *testing.T) {
	// TODO change on integration tests
	t.Skip("test works with preconditions")

	cl := makeTestClient(t)

	ts := strconv.Itoa(0)

	req := provider.ConnectRequest{
		UserID:     crypto.Hash(pbK),
		ConsumerID: "f65af5d64000c7cd2883f4910eb69086f9d6e6635c744e62afcfab58b938ee25",
		Auth: &consumer.Auth{
			CreationDate:    ts,
			SignatureScheme: "bls0chain",
			Signature:       sign,
			PublicKey:       pbK,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	_, err := cl.Connect(ctx, &req)
	require.NoError(t, err)
}

//nolint:unused
func makeTestClient(t *testing.T) provider.AccessPointClient {
	conn, err := grpc.Dial(":7041", grpc.WithInsecure())
	require.NoError(t, err)

	return provider.NewAccessPointClient(conn)
}
