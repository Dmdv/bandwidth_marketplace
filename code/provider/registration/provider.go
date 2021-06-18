package registration

import (
	"encoding/json"
	"time"

	"github.com/0chain/bandwidth_marketplace/code/core/errors"
	"github.com/0chain/bandwidth_marketplace/code/core/log"
	"github.com/0chain/bandwidth_marketplace/code/core/node"
	"github.com/0chain/bandwidth_marketplace/code/core/provider"
	"github.com/0chain/bandwidth_marketplace/code/core/transaction"
	"github.com/0chain/gosdk/zcncore"
	"go.uber.org/zap"

	"provider/config"
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
	zcncore.SetLogLevel(logLevelFromStr(cfg.Logging.Level))
	if err := zcncore.InitZCNSDK(cfg.ServerChain.BlockWorker, cfg.ServerChain.SignatureScheme); err != nil {
		errors.ExitErr("error while init zcn sdk", err, 2)
	}

	if err := zcncore.SetWalletInfo(wallet, false); err != nil {
		errors.ExitErr("error while init wallet", err, 2)
	}
}

// RegisterWithRetries registers the current self node in blockchain with retries.
// Checks self node for already registered provider node with same ID,
// and if not executes transaction that registers provider and verifies including transaction in blockchain.
//
// If an error occurs during execution, the program terminates with code 2 and the last error will be written in os.Stderr.
//
// RegisterWithRetries should be used only once while application is starting.
func RegisterWithRetries(numTries int, cfg config.Terms) {
	log.Logger.Info("Start Provider registration ...")

	var (
		ind        int
		registered bool
		txnHash    string
		err        error
	)
	for ind < numTries && !registered {
		registered, _ = isProviderRegistered()
		if registered {
			log.Logger.Info("Provider already registered in blockchain")
			break
		}

		txnHash, err = register(cfg)
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
		errors.ExitErr("error while registering provider", err, 2)
	}

	log.Logger.Info("Node is registered as Magma Provider in blockchain", zap.String("ID", node.GetSelfNode().ID()))
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

// register registers client's wallet to miners, and executes MagmaSC transaction.RegisterProviderFuncName function.
func register(cfg config.Terms) (string, error) {
	err := zcncore.RegisterToMiners(node.GetSelfNode().GetWallet(), &walletCallback{})
	if err != nil {
		return "", err
	}

	txn, err := transaction.NewTransactionEntity()
	if err != nil {
		return "", err
	}

	prov := provider.Provider{
		ID: node.GetSelfNode().ID(),
		Terms: provider.Terms{
			Price: cfg.Price,
			QoS: provider.QoS{
				DownloadMBPS: cfg.QoS.DownloadMBPS,
				UploadMBPS:   cfg.QoS.UploadMBPS,
			},
		},
	}
	input, err := json.Marshal(prov)
	if err != nil {
		return "", err
	}

	return txn.ExecuteSmartContract(transaction.MagmaSCAddress, transaction.RegisterProviderFuncName, string(input), 0)
}

// isProviderRegistered makes request to MagmaSC transaction.GetAllProviders rest point and looks for self node
// in registered providers.
func isProviderRegistered() (bool, error) {
	resp, err := transaction.MakeSCRestAPICall(transaction.MagmaSCAddress, transaction.GetAllProvidersRP, nil)
	if err != nil {
		return false, err
	}

	providers := make([]provider.Provider, 0)
	if err := json.Unmarshal(resp, &providers); err != nil {
		return false, err
	}

	for _, p := range providers {
		if node.GetSelfNode().ID() == p.ID {
			return true, nil
		}
	}
	return false, nil
}
