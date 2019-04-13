package cryptography

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"

	"golang.org/x/crypto/nacl/box"
)

const (
	// NonceSize is a size of nonce used in Nacl in bytes
	NonceSize      = 24
	privateKeyName = "private.key"
	publicKeyName  = "public.key"
)

// Box is structure performing cryptographic operations
// using asymmetric cryptographic algorithm defined in:
// https://nacl.cr.yp.to
type Box struct {
	privateKey Key32
	publicKey  Key32
}

// NewBox creates new instance of the Box
func NewBox(privateKey Key32, publicKey Key32) *Box {
	return &Box{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// NewOneShotBox creates Box instance with randomly generated keys
// It's supposed to be one time use.
func NewOneShotBox() (*Box, error) {
	publicKey, privateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Box{
		publicKey:  *publicKey,
		privateKey: *privateKey,
	}, nil
}

func storeKey(key *[32]byte, filepath string) error {
	return ioutil.WriteFile(filepath, key[:], 0600)
}

// Encrypt encrypts message using recipient public key.
// If encryption fails proper error is returned.
// Node: the message is also signed with senders private key.
func (n *Box) Encrypt(message []byte, recipientPublicKey *Key32) ([]byte, error) {
	var (
		nonce [24]byte
	)

	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}

	out := box.Seal(nonce[:], message, &nonce, (*[32]byte)(recipientPublicKey), (*[32]byte)(&n.privateKey))
	return out, nil
}

// Decrypt decrypts message from sender.
// For this operation sender public key is needed.
func (n *Box) Decrypt(message []byte, senderPublicKey *Key32) ([]byte, error) {
	var decryptNonce [24]byte
	copy(decryptNonce[:], message[:24])

	decrypted, ok := box.Open(nil, message[24:], &decryptNonce, (*[32]byte)(senderPublicKey), (*[32]byte)(&n.privateKey))
	if !ok {
		return nil, fmt.Errorf("failed to decrypt message")
	}
	return decrypted, nil
}

// GetPublicKey returns public key
func (n *Box) GetPublicKey() *Key32 {
	return &n.publicKey
}

// GetPrivateKey returns private key
func (n *Box) GetPrivateKey() *Key32 {
	return &n.privateKey
}

// recompute calculates the shared key between peersPublicKey and privateKey and returns it
// The shared key can be used with EncryptAfterPrecomputation and DecryptAfterPrecomputation to speed up processing
// hen using the same pair of keys repeatedly.
func (n *Box) Precompute(recipientPublicKey *Key32) (sharedKey Key32) {
	box.Precompute((*[32]byte)(&sharedKey), (*[32]byte)(recipientPublicKey), (*[32]byte)(&n.privateKey))
	return sharedKey
}

func (n *Box) EncryptAfterPrecomputation(message []byte, sharedKey *Key32) ([]byte, error) {
	var (
		nonce [24]byte
	)

	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}

	out := box.SealAfterPrecomputation(nonce[:], message, &nonce, (*[32]byte)(sharedKey))
	return out, nil
}

// Decrypt decrypts message from sender.
// For this operation sender public key is needed.
func (n *Box) DecryptAfterPrecomputation(message []byte, sharedKey *Key32) ([]byte, error) {
	var decryptNonce [24]byte
	copy(decryptNonce[:], message[:24])

	decrypted, ok := box.OpenAfterPrecomputation(nil, message[24:], &decryptNonce, (*[32]byte)(sharedKey))
	if !ok {
		return nil, fmt.Errorf("failed to decrypt message")
	}
	return decrypted, nil
}

// BoxEncrypt encrypts message using recipient public key.
// If encryption fails proper error is returned.
// Node: the message is also signed with senders private key.
func BoxEncrypt(message []byte, recipientPublicKey *Key32, senderPrivateKey *Key32) ([]byte, error) {
	var (
		nonce [24]byte
	)

	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}

	out := box.Seal(nonce[:], message, &nonce, (*[32]byte)(recipientPublicKey), (*[32]byte)(senderPrivateKey))
	return out, nil
}

// BoxDecrypt decrypts message from sender.
// For this operation sender public key is needed.
func BoxDecrypt(message []byte, recipientPrivateKey *Key32, senderPublicKey *Key32) ([]byte, error) {
	var decryptNonce [24]byte
	copy(decryptNonce[:], message[:24])

	decrypted, ok := box.Open(nil, message[24:], &decryptNonce, (*[32]byte)(senderPublicKey), (*[32]byte)(recipientPrivateKey))
	if !ok {
		return nil, fmt.Errorf("failed to decrypt message")
	}
	return decrypted, nil
}
