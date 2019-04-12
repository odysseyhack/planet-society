package cryptography

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/nacl/secretbox"
)

// SecretBox is structure performing cryptographic operations
// using symmetric cryptographic algorithm defined in:
// https://nacl.cr.yp.to
type SecretBox struct {
	key Key32
}

// NewSecretBox creates new instance of SecretBox
func NewSecretBox(key Key32) *SecretBox {
	return &SecretBox{key: key}
}

// Encrypt encrypts message using symmetric algorithm
func (s *SecretBox) Encrypt(message []byte) ([]byte, error) {
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, err
	}

	return secretbox.Seal(nonce[:], message, &nonce, (*[32]byte)(&s.key)), nil
}

// Decrypt decrypts message using symmetric algorithm
func (s *SecretBox) Decrypt(encrypted []byte) ([]byte, error) {
	var decryptNonce [24]byte
	copy(decryptNonce[:], encrypted[:24])

	decrypted, ok := secretbox.Open(nil, encrypted[24:], &decryptNonce, (*[32]byte)(&s.key))
	if !ok {
		return nil, fmt.Errorf("failed to decrypt message")
	}

	return decrypted, nil
}
