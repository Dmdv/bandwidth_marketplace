package server

import (
	"context"
	"strconv"
	"testing"

	"github.com/0chain/gosdk/core/zcncrypto"
	"github.com/MurashovVen/bandwidth-marketplace/code/core/crypto"
	"github.com/MurashovVen/bandwidth-marketplace/code/pb/consumer"
	"github.com/MurashovVen/bandwidth-marketplace/code/pb/magma"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const (
	prK = "0deb9bc79c4061252ab671109ec1f572bd495d9114c5b843d29f49c4648b5b0c"
	pbK = "3d8b2c47542404b0059a4c320dca873645d0a3259532a8e17e4cd10db2ef2012b0ba6ef46c5dffab0bb43c05cfc3da7004adf5fe5497e339b298f742fa9843a3"
)

func Test_magmaServer_Connect(t *testing.T) {
	cl := makeTestClient(t)

	scheme := zcncrypto.NewBLS0ChainScheme()
	scheme.PrivateKey = prK
	scheme.PublicKey = pbK

	ts := strconv.Itoa(int(0))
	sign, err := scheme.Sign(crypto.Hash(ts))
	require.NoError(t, err)

	req := magma.ConnectRequest{
		UserID:     pbK,
		ProviderID: "7a90e6790bcd3d78422d7a230390edc102870fe58c15472073922024985b1c7d",
		Auth: &consumer.Auth{
			CreationDate:    ts,
			SignatureScheme: "bls0chain",
			Signature:       sign,
		},
	}
	_, err = cl.Connect(context.Background(), &req)
	require.NoError(t, err)

	//fmt.Println(resp.String())
}

func makeTestClient(t *testing.T) magma.MagmaClient {
	conn, err := grpc.Dial(":8081", grpc.WithInsecure())
	require.NoError(t, err)

	return magma.NewMagmaClient(conn)
}
