package registration

import (
	"encoding/json"
	"time"

	"github.com/0chain/bandwidth_marketplace/code/core/consumer"
	"github.com/0chain/bandwidth_marketplace/code/core/errors"
	"github.com/0chain/bandwidth_marketplace/code/core/log"
	"github.com/0chain/bandwidth_marketplace/code/core/node"
	"github.com/0chain/bandwidth_marketplace/code/core/transaction"
	"github.com/0chain/gosdk/zcncore"
	"go.uber.org/zap"

	"consumer/config"
)

// Setup runs zcncore.SetLogFile, zcncore.SetLogLevel using provided log directory,
// zcncore.InitZCNSDK using provided config.Config, and zcncore.SetWalletInfo with provided wallet.
//
// If an error occurs during execution, the program terminates with code 2 and the error will be written in os.Stderr.
//
// Setup should be used only once while application is starting.
func Setup(cfg *config.Config, wallet string) {
	var logName = cfg.CLConfig.LogDir + "/zsdk.log"
	zcncore.SetLogFile(logName, false)
	zcncore.SetLogLevel(intLevelFromStr(cfg.Logging.Level))
	if err := zcncore.InitZCNSDK(cfg.ServerChain.BlockWorker, cfg.ServerChain.SignatureScheme); err != nil {
		errors.ExitErr("error while init zcn sdk", err, 2)
	}

	if err := zcncore.SetWalletInfo(wallet, false); err != nil {
		errors.ExitErr("error while init wallet", err, 2)
	}
}

// RegisterWithRetries registers the current self node in blockchain with retries.
// Checks self node for already registered bandwidth-marketplace with same ID,
// and if not executes transaction.RegisterConsumerFuncName and verifies including transaction in blockchain.
//
// If an error occurs during execution, the program terminates with code 2 and the last error will be written in os.Stderr.
//
// RegisterWithRetries should be used only once while application is starting.
func RegisterWithRetries(numTries int) {
	var (
		ind        int
		registered bool
		txnHash    string
		err        error
	)
	for ind < numTries && !registered {
		registered, _ = isConsumerRegistered()
		if registered {
			log.Logger.Info("Consumer is already registered in blockchain")
			break
		}

		txnHash, err = register()
		if err != nil {
			ind++
			time.Sleep(time.Second / 2)
			continue
		}

		_, err := transaction.VerifyTransaction(txnHash)
		if err != nil {
			ind++
			log.Logger.Error("Verifying transaction failed. Sleeping for 5 seconds ...", zap.Error(err))
			time.Sleep(time.Second * 5)
			continue
		}

		break
	}

	if err != nil {
		errors.ExitErr("error while registering consumer: %v\n", err, 2)
	}

	log.Logger.Info("Node is registered as Magma Consumer in blockchain", zap.String("ID", node.GetSelfNode().ID()))
}

// walletCallback provides callback struct for operations with wallet.
type walletCallback struct{}

// OnWalletCreateComplete implements zcncore.WalletCallback interface.
func (wb *walletCallback) OnWalletCreateComplete(status int, wallet string, err string) {
	if err == "" {
		err = "no error"
	}
	log.Logger.Info("Creating wallet completed", zap.Int("status", status), zap.String("wallet", wallet), zap.String("error", err))
}

// register registers client's wallet to miners, and executes MagmaSC transaction.RegisterConsumerFuncName function.
func register() (string, error) {
	err := zcncore.RegisterToMiners(node.GetSelfNode().GetWallet(), &walletCallback{})
	if err != nil {
		return "", err
	}

	txn, err := transaction.NewTransactionEntity()
	if err != nil {
		return "", err
	}

	return txn.ExecuteSmartContract(transaction.MagmaSCAddress, transaction.RegisterConsumerFuncName, "", 0)
}

// isConsumerRegistered makes request to MagmaSC transaction.GetAllConsumersRP rest point and looks for self node
// in registered consumers.
func isConsumerRegistered() (bool, error) {
	resp, err := transaction.MakeSCRestAPICall(transaction.MagmaSCAddress, transaction.GetAllConsumersRP, nil)
	if err != nil {
		return false, err
	}

	consumers := make([]consumer.Consumer, 0)
	if err := json.Unmarshal(resp, &consumers); err != nil {
		return false, err
	}

	for _, c := range consumers {
		if node.GetSelfNode().ID() == c.ID {
			return true, nil
		}
	}
	return false, nil
}
