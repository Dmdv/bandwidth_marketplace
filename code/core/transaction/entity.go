package transaction

import (
	"encoding/json"
	"sync"

	"github.com/0chain/gosdk/core/util"
	"github.com/0chain/gosdk/zcncore"
	"go.uber.org/zap"

	"github.com/0chain/bandwidth_marketplace/code/core/chain"
	"github.com/0chain/bandwidth_marketplace/code/core/errors"
	"github.com/0chain/bandwidth_marketplace/code/core/log"
	"github.com/0chain/bandwidth_marketplace/code/core/node"
	"github.com/0chain/bandwidth_marketplace/code/core/time"
)

type (
	// Transaction entity that encapsulates the transaction related data and meta data.
	Transaction struct {
		Hash              string         `json:"hash,omitempty"`
		Version           string         `json:"version,omitempty"`
		ClientID          string         `json:"client_id,omitempty"`
		PublicKey         string         `json:"public_key,omitempty"`
		ToClientID        string         `json:"to_client_id,omitempty"`
		ChainID           string         `json:"chain_id,omitempty"`
		TransactionData   string         `json:"transaction_data,omitempty"`
		Value             int64          `json:"transaction_value,omitempty"`
		Signature         string         `json:"signature,omitempty"`
		CreationDate      time.Timestamp `json:"creation_date,omitempty"`
		TransactionType   int            `json:"transaction_type,omitempty"`
		TransactionOutput string         `json:"transaction_output,omitempty"`
		OutputHash        string         `json:"txn_output_hash"`

		scheme zcncore.TransactionScheme
		wg     *sync.WaitGroup
	}
)

// NewTransactionEntity creates Transaction with initialized fields.
// Sets version, client ID, creation date, public key and creates internal zcncore.TransactionScheme.
func NewTransactionEntity() (*Transaction, error) {
	selfNode := node.GetSelfNode()

	txn := &Transaction{
		Version:      "1.0",
		ClientID:     selfNode.ID(),
		CreationDate: time.Now(),
		ChainID:      chain.GetServerChain().ID,
		wg:           new(sync.WaitGroup),
		PublicKey:    selfNode.PublicKey(),
	}
	zcntxn, err := zcncore.NewTransaction(txn, 0)
	if err != nil {
		return nil, err
	}
	txn.scheme = zcntxn

	return txn, nil
}

// ExecuteSmartContract executes function of smart contract with provided address.
//
// Returns hash of executed transaction.
func (t *Transaction) ExecuteSmartContract(address, funcName, input string, val int64) (string, error) {
	t.wg.Add(1)
	err := t.scheme.ExecuteSmartContract(address, funcName, input, val)
	if err != nil {
		t.wg.Done()
		return "", err
	}
	t.wg.Wait()

	t.Hash = t.scheme.GetTransactionHash()
	if len(t.scheme.GetTransactionError()) > 0 {
		return "", errors.NewError("transaction_send_error", t.scheme.GetTransactionError())
	}

	return t.Hash, nil
}

type (
	verifyOutput struct {
		Txn          *Transaction `json:"txn"`
		Confirmation confirmation `json:"confirmation"`
	}

	// confirmation represents the acceptance that a transaction is included into the block chain.
	confirmation struct {
		Version               string         `json:"version"`
		Hash                  string         `json:"hash"`
		BlockHash             string         `json:"block_hash"`
		PreviousBlockHash     string         `json:"previous_block_hash"`
		Transaction           *Transaction   `json:"txn,omitempty"`
		CreationDate          time.Timestamp `json:"creation_date"`
		MinerID               string         `json:"miner_id"`
		Round                 int64          `json:"round"`
		Status                int            `json:"transaction_status" msgpack:"sot"`
		RoundRandomSeed       int64          `json:"round_random_seed"`
		MerkleTreeRoot        string         `json:"merkle_tree_root"`
		MerkleTreePath        *util.MTPath   `json:"merkle_tree_path"`
		ReceiptMerkleTreeRoot string         `json:"receipt_merkle_tree_root"`
		ReceiptMerkleTreePath *util.MTPath   `json:"receipt_merkle_tree_path"`
	}
)

// Verify checks including of transaction in the blockchain.
func (t *Transaction) Verify() error {
	if err := t.scheme.SetTransactionHash(t.Hash); err != nil {
		return err
	}

	t.wg.Add(1)
	err := t.scheme.Verify() // TODO need to be refactored, change wait group on channel
	if err != nil {
		t.wg.Done()
		return err
	}
	t.wg.Wait()
	if len(t.scheme.GetVerifyError()) > 0 {
		return errors.NewError("transaction_verify_error", t.scheme.GetVerifyError())
	}

	// just checking, can be removed
	vo := new(verifyOutput)
	if err := json.Unmarshal([]byte(t.scheme.GetVerifyOutput()), vo); err != nil {
		return err
	}

	return nil
}

// OnTransactionComplete implements zcncore.TransactionCallback interface.
func (t *Transaction) OnTransactionComplete(zcnTxn *zcncore.Transaction, status int) {
	t.wg.Done()

	err := zcnTxn.GetTransactionError()
	if err == "" {
		err = "no error"
	}

	log.Logger.Info("Transaction completed",
		zap.String("status", TxnStatus(status).String()),
		zap.String("hash", zcnTxn.GetTransactionHash()),
		zap.String("error", err),
	)
}

// OnVerifyComplete implements zcncore.TransactionCallback interface.
func (t *Transaction) OnVerifyComplete(zcnTxn *zcncore.Transaction, status int) {
	t.wg.Done()

	err := zcnTxn.GetVerifyError()
	if err == "" {
		err = "no error"
	}

	log.Logger.Info("Transaction verified",
		zap.String("status", TxnStatus(status).String()),
		zap.String("hash", zcnTxn.GetTransactionHash()),
		zap.String("error", err),
	)
}

// OnAuthComplete implements zcncore.TransactionCallback interface.
func (t *Transaction) OnAuthComplete(zcnTxn *zcncore.Transaction, status int) {
	log.Logger.Info("Transaction authenticated",
		zap.String("status", TxnStatus(status).String()),
		zap.String("hash", zcnTxn.GetTransactionHash()),
	)
}
