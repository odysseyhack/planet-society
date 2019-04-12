package cryptography

import (
	"golang.org/x/crypto/blake2b"
)

// Hash16 creates 16 byte hash from slice of bytes.
// The returned id is created with blake2b hashing function
// and trimmed to fit in size.
func Hash16(data []byte) (ret [16]byte) {
	hash := blake2b.Sum256(data)
	copy(ret[:], hash[:])
	return ret
}
