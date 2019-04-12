package cryptography

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/nacl/sign"
)

// Signer signs and verifies small messages using public-key cryptography.
type Signer struct {
	privateKey Key64
	publicKey  Key32
}

// Sign signs message using public key cryptography
func (s *Signer) Sign(message []byte) ([]byte, error) {
	if message == nil {
		return nil, fmt.Errorf("sign: message is nil")
	}

	return sign.Sign(nil, message, (*[64]byte)(&s.privateKey)), nil
}

// Verify verifies message using public key cryptography
func (s *Signer) Verify(message []byte) ([]byte, error) {
	if message == nil {
		return nil, fmt.Errorf("verify: message is nil")
	}

	message, ok := sign.Open(nil, message, (*[32]byte)(&s.publicKey))
	if !ok {
		return nil, fmt.Errorf("verify: verification failed")
	}

	return message, nil
}

// GenerateSigner generates set of cryptographic keys used for signing
func GenerateSigner() (*Signer, error) {
	publicKey, privateKey, err := sign.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	signer := &Signer{}
	copy(signer.publicKey[:], publicKey[:])
	copy(signer.privateKey[:], privateKey[:])
	return signer, nil
}
