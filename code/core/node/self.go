package node

import (
	"encoding/hex"
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	"github.com/0chain/gosdk/core/zcncrypto"

	"github.com/0chain/bandwidth_marketplace/code/core/crypto"
	"github.com/0chain/bandwidth_marketplace/code/core/errors"
)

// Node represent self node type.
type Node struct {
	url       string
	wallet    *zcncrypto.Wallet
	id        string
	publicKey string
	startTime time.Time
}

var (
	self = &Node{}
)

// Start writes to self node current time.
//
// Start should be used only once while application is starting.
func Start() {
	self.startTime = time.Now()
}

// GetSelfNode returns self Node variable described as package variable.
func GetSelfNode() *Node {
	return self
}

// SetKeys creates wallet and sets provided keys to Node.
//
// If an error occurs during execution, the program terminates with code 2 and the error will be written in os.Stderr.
//
// SetKeys should be used only once while application is starting.
func (sn *Node) SetKeys(publicKey, privateKey string) {
	publicKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		errors.ExitErr("error while hex-decode public key", err, 2)
	}

	sn.wallet = &zcncrypto.Wallet{
		ClientID:  crypto.Hash(publicKeyBytes),
		ClientKey: publicKey,
		Keys:      make([]zcncrypto.KeyPair, 1),
		Version:   zcncrypto.CryptoVersion,
	}
	sn.wallet.Keys[0].PublicKey = publicKey
	sn.wallet.Keys[0].PrivateKey = privateKey

	sn.id = sn.wallet.ClientID
	sn.publicKey = publicKey
}

// SetHostURL is a setter for nodes URL.
//
// Should be used only once while application is starting.
func (sn *Node) SetHostURL(host string, port int) {
	if host == "" {
		host = "localhost"
	}

	uri := url.URL{
		Scheme: "http",
		Host:   host + ":" + strconv.Itoa(port),
	}

	sn.url = uri.String()
}

// Validate returns error if Node.ID is empty and nil error if it is not.
func (sn *Node) Validate() error {
	if sn.id == "" {
		return errors.NewError("invalid_id", "empty node id")
	}
	return nil
}

// GetURLBase returns nodes URL.
func (sn *Node) GetURLBase() string {
	return sn.url
}

// GetWalletString returns marshaled to JSON string nodes wallet.
func (sn *Node) GetWalletString() string {
	walletStr, _ := json.Marshal(sn.wallet)
	return string(walletStr)
}

// ID returns id of Node.
func (sn *Node) ID() string {
	return sn.id
}

// PublicKey returns id of Node.
func (sn *Node) PublicKey() string {
	return sn.publicKey
}

// StartTime returns time when Node is started.
func (sn *Node) StartTime() time.Time {
	return sn.startTime
}

func (sn *Node) GetWallet() *zcncrypto.Wallet {
	return sn.wallet
}
