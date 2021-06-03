package transaction

const (
	// MagmaSCAddress represents the address of the Magma smart contract.
	// Used while making requests to smart contract's rest points and executing smart contracts functions.
	MagmaSCAddress = "11f8411db41e34cea7c100f19faff32da8f3cd5a80635731cec06f32d08089be"

	// GetAllConsumersRP represents MagmaSC relative path.
	// Used to list all registered in the blockchain rest points.
	GetAllConsumersRP = "/getAllConsumers"

	// RegisterConsumerFuncName represents MagmaSC relative path.
	// Used to register bandwidth-marketplace's node.
	RegisterConsumerFuncName = "register_consumer"

	// AcceptTermsFuncName represents MagmaSC function.
	// Used to confirm bandwidth-marketplace's acceptance of provider service terms.
	AcceptTermsFuncName = "accept_terms"

	// RegisterProviderFuncName represents MagmaSC relative path.
	// Used to register bandwidth-marketplace's node.
	RegisterProviderFuncName = "register_provider"

	// GetAllProvidersRP represents MagmaSC relative path.
	// Used to list all registered in the blockchain rest points.
	GetAllProvidersRP = "/getAllProviders"
)

type (
	// TxnStatus represented zcncore.TransactionCallback operations statuses.
	TxnStatus int
)

const (
	// StatusSuccess represent zcncore.StatusSuccess.
	StatusSuccess TxnStatus = iota
	// StatusNetworkError represent zcncore.StatusNetworkError.
	StatusNetworkError
	// StatusError represent zcncore.StatusError.
	StatusError
	// StatusRejectedByUser represent zcncore.StatusRejectedByUser.
	StatusRejectedByUser
	// StatusInvalidSignature represent zcncore.StatusInvalidSignature.
	StatusInvalidSignature
	// StatusAuthError represent zcncore.StatusAuthError.
	StatusAuthError
	// StatusAuthVerifyFailed represent zcncore.StatusAuthVerifyFailed.
	StatusAuthVerifyFailed
	// StatusAuthTimeout represent zcncore.StatusAuthTimeout.
	StatusAuthTimeout
	// StatusUnknown represent zcncore.StatusUnknown.
	StatusUnknown = -1
)

// String returns represented in string format TxnStatus.
func (ts TxnStatus) String() string {
	switch ts {
	case StatusSuccess:
		return "success"

	case StatusNetworkError:
		return "network error"

	case StatusError:
		return "error"

	case StatusRejectedByUser:
		return "rejected byt user"

	case StatusInvalidSignature:
		return "invalid signature"

	case StatusAuthError:
		return "auth error"

	case StatusAuthVerifyFailed:
		return "auth verify error"

	case StatusAuthTimeout:
		return "auth timeout error"

	case StatusUnknown:
		return "unknown"

	default:
		return ""
	}
}
