package sha256

import (
	"crypto/sha256"
)

func Hash(bytes []byte) []byte {
	hash := sha256.New()
	hash.Write(bytes)
	return hash.Sum(nil)
}
