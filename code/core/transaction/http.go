package transaction

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/0chain/gosdk/core/util"
	"github.com/0chain/gosdk/zcncore"
	"go.uber.org/zap"

	"github.com/0chain/bandwidth_marketplace/code/core/errors"
	chttp "github.com/0chain/bandwidth_marketplace/code/core/http"
	"github.com/0chain/bandwidth_marketplace/code/core/log"
)

const (
	// ScRestApiUrl represents base URL path to execute smart contract rest points.
	ScRestApiUrl = "v1/screst/"
)

// VerifyTransaction verifies including in blockchain transaction with provided hash.
//
// If execution completed with no error, returns Transaction with provided hash.
func VerifyTransaction(txnHash string) (*Transaction, error) {
	txn, err := NewTransactionEntity()
	if err != nil {
		return nil, err
	}

	txn.Hash = txnHash
	err = txn.Verify()
	if err != nil {
		return nil, err
	}
	return txn, nil
}

// MakeSCRestAPICall calls smart contract with provided address
// and makes retryable request to smart contract resource with provided relative path using params.
func MakeSCRestAPICall(scAddress string, relativePath string, params map[string]string) ([]byte, error) {
	var (
		resMaxCounterBody []byte

		hashMaxCounter int
		hashCounters   = make(map[string]int)

		sharders             = extractSharders()
		remainingShardersNum = len(sharders)
	)

	for _, sharder := range sharders {
		// Make one or more requests (in case of unavailability, see 503/504 errors)
		var (
			err     error
			resp    *http.Response
			retries = 3

			netClient = chttp.NewClient()
			u         = makeScURL(params, sharder, scAddress, relativePath)
		)
		for retries > 0 {
			resp, err = netClient.Get(u.String())
			if err != nil || resp.StatusCode != http.StatusServiceUnavailable && resp.StatusCode != http.StatusGatewayTimeout {
				break
			}
			_ = resp.Body.Close()
			retries--
		}
		if err != nil {
			log.Logger.Error("Error getting response for sc rest api", zap.Error(err), zap.Any("sharder_url", sharder))
			remainingShardersNum--
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resBody, _ := ioutil.ReadAll(resp.Body)
			log.Logger.Error("Got error response from sc rest api", zap.Any("response", string(resBody)))
			_ = resp.Body.Close()
			continue
		}

		hash, resBody, err := hashAndBytesOfReader(resp.Body)
		if err != nil {
			log.Logger.Error("Error reading response", zap.Error(err))
			_ = resp.Body.Close()
			continue
		}

		hashCounters[hash]++
		if hashCounters[hash] > hashMaxCounter {
			hashMaxCounter = hashCounters[hash]
			resMaxCounterBody = resBody
		}

		_ = resp.Body.Close()
	}

	var err error

	// is it less than or equal to 50%
	if hashMaxCounter <= (remainingShardersNum / 2) {
		err = errors.NewError("invalid_response", "can not make request to sharders")
	} else {
		return resMaxCounterBody, nil
	}

	return nil, err
}

// hashAndBytesOfReader computes hash of readers data and returns hash encoded to hex and bytes of reader data.
// If error occurs while reading data from reader, it returns non nil error.
func hashAndBytesOfReader(r io.Reader) (hash string, reader []byte, err error) {
	h := sha1.New()
	teeReader := io.TeeReader(r, h)
	readerBytes, err := ioutil.ReadAll(teeReader)
	if err != nil {
		return "", nil, err
	}

	return hex.EncodeToString(h.Sum(nil)), readerBytes, nil
}

// extractSharders returns string slice of randomly ordered sharders existing in the current network.
func extractSharders() []string {
	network := zcncore.GetNetwork()
	return util.GetRandom(network.Sharders, len(network.Sharders))
}

// makeScURL creates url.URL to make smart contract request to sharder.
func makeScURL(params map[string]string, sharder, scAddress, relativePath string) *url.URL {
	uString := fmt.Sprintf("%v/%v/%v%v", sharder, ScRestApiUrl, scAddress, relativePath)
	u, _ := url.Parse(uString)
	q := u.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()

	return u
}
