package provider

import (
	"encoding/json"

	"github.com/0chain/bandwidth_marketplace/code/core/node"
	"github.com/0chain/bandwidth_marketplace/code/core/time"
	"github.com/0chain/bandwidth_marketplace/code/core/transaction"
)

type (
	// Acknowledgment represents accepting Provider's terms by Consumer.
	Acknowledgment struct {
		AccessPointID string        `json:"access_point_id"`
		ConsumerID    string        `json:"consumer_id"`
		ProviderID    string        `json:"provider_id"`
		SessionID     string        `json:"session_id"`
		ProviderTerms ProviderTerms `json:"provider_terms"`
	}

	// Provider represents providers node stored in block chain.
	Provider struct {
		ID    string        `json:"id"`
		Terms ProviderTerms `json:"terms"`
	}

	// ProviderTerms represents information of provider and services terms.
	ProviderTerms struct {
		Terms
		QoS QoS `json:"qos"`
	}

	// Terms represents information of provider terms.
	Terms struct {
		Price     int64          `json:"price"`      // per byte
		Volume    int64          `json:"volume"`     // in bytes
		ExpiredAt time.Timestamp `json:"expired_at"` // valid till
	}

	// QoS represents a Quality of Service and contains uploading and downloading speed
	// represented in megabits per second.
	QoS struct {
		DownloadMBPS int64 `json:"download_mbps"` // megabits per second
		UploadMBPS   int64 `json:"upload_mbps"`   // megabits per second
	}
)

// ExecuteAcceptTerms creates Acknowledgment with provided args, executes MagmaSC transaction.AcceptTermsFuncName and
// verifies including the transaction in the blockchain.
//
// Returns resulted Acknowledgment, which ID equals executed transaction hash.
func ExecuteAcceptTerms(provID, apID, sessID string) (*Acknowledgment, error) {
	txn, err := transaction.NewTransactionEntity()
	if err != nil {
		return nil, err
	}

	ackn := Acknowledgment{
		ProviderID:    provID,
		AccessPointID: apID,
		SessionID:     sessID,
	}
	input, err := json.Marshal(&ackn)
	if err != nil {
		return nil, err
	}
	txnHash, err := txn.ExecuteSmartContract(transaction.MagmaSCAddress, transaction.AcceptTermsFuncName, string(input), 0)
	if err != nil {
		return nil, err
	}

	if _, err := transaction.VerifyTransaction(txnHash); err != nil {
		return nil, err
	}

	ackn.ConsumerID = node.GetSelfNode().ID()
	return &ackn, err
}
