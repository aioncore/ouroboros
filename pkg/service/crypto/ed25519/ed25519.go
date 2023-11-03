package ed25519

import (
	"crypto/ed25519"
	"crypto/rand"
	"github.com/aioncore/ouroboros/pkg/service/crypto/sha256"
	"io"
)

// GeneratePrivateKey generates a new ed25519 private key.
// It uses OS randomness in conjunction with the current global random seed
// in tendermint/libs/common to generate the private key.
func GeneratePrivateKey() (ed25519.PrivateKey, error) {
	return generatePrivateKey(rand.Reader)
}

// genPrivateKey generates a new ed25519 private key using the provided reader.
func generatePrivateKey(rand io.Reader) (ed25519.PrivateKey, error) {
	seed := make([]byte, 32)

	_, err := io.ReadFull(rand, seed)
	if err != nil {
		return nil, err
	}

	return ed25519.NewKeyFromSeed(seed), nil
}

// GeneratePrivateKeyFromSecret hashes the secret with SHA2, and uses
// that 32 byte output to create the private key.
// NOTE: secret should be the output of a KDF like bcrypt,
// if it's derived from user input.
func GeneratePrivateKeyFromSecret(secret []byte) ed25519.PrivateKey {
	seed := sha256.Hash(secret)
	return ed25519.NewKeyFromSeed(seed)
}
