package crypto

import (
	"bufio"
	"io"
	"os"

	"github.com/0chain/gosdk/core/common"
	"github.com/0chain/gosdk/core/zcncrypto"

	"github.com/0chain/bandwidth_marketplace/code/core/errors"
)

// ReadKeysFile reads file existing in keysFile dir and parses public and private keys from file.
//
// If the reads succeeds, it returns public and private keys, either if an error occurs during execution,
// the program terminates with code 2 and the error will be written in os.Stderr.
//
// ReadKeysFile should be used only once while application is starting.
func ReadKeysFile(keysFile string) (publicKey string, privateKey string) {
	reader, err := os.Open(keysFile)
	if err != nil {
		errors.ExitErr("error while open keys file", err, 2)
	}

	publicKey, privateKey = readKeys(reader)
	err = reader.Close()
	if err != nil {
		errors.ExitErr("error while close keys file", err, 2)
	}

	return publicKey, privateKey
}

// readKeys reads a publicKey and a privateKey from a io.Reader passed in args.
// They are assumed to be in two separate lines one followed by the other.
func readKeys(reader io.Reader) (publicKey string, privateKey string) {
	scanner := bufio.NewScanner(reader)
	scanner.Scan()
	publicKey = scanner.Text()
	scanner.Scan()
	privateKey = scanner.Text()
	scanner.Scan()

	return publicKey, privateKey
}

// Verify verifies passed signature of the passed hash with passed public key using the signature scheme.
func Verify(publicKey, signature, hash, scheme string) (bool, error) {
	signScheme := zcncrypto.NewSignatureScheme(scheme)
	if signScheme != nil {
		err := signScheme.SetPublicKey(publicKey)
		if err != nil {
			return false, err
		}
		return signScheme.Verify(signature, hash)
	}
	return false, common.NewError("invalid_signature_scheme", "invalid signature scheme")
}
