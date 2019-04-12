package cryptography

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"

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

// GenerateKeysToFile generates pair of public and private keys.
// It stores keys in directory defined by path under names:
// public.key and private.key
// Keys will not be generate if files already exist.
func GenerateKeysToFile(path string) error {
	var (
		privateKeyPath = fmt.Sprintf("%s/%s", path, privateKeyName)
		publicKeyPath  = fmt.Sprintf("%s/%s", path, publicKeyName)
	)

	if _, err := os.Stat(privateKeyPath); err == nil {
		return fmt.Errorf("private key already exist")
	}

	if _, err := os.Stat(publicKeyPath); err == nil {
		return fmt.Errorf("public key already exist")
	}

	publicKey, privateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(path, 0700); err != nil {
		return err
	}

	if err := storeKey(publicKey, publicKeyPath); err != nil {
		return err
	}

	if err := storeKey(privateKey, privateKeyPath); err != nil {
		return err
	}

	return nil
}

func storeKey(key *[32]byte, filepath string) error {
	return ioutil.WriteFile(filepath, key[:], 0600)
}

// BoxLoad loads public and private keys from directory
// defined in path and stores them in NaCl structure
func BoxLoad(path string) (*Box, error) {
	var (
		privateKeyPath = fmt.Sprintf("%s/%s", path, privateKeyName)
		publicKeyPath  = fmt.Sprintf("%s/%s", path, publicKeyName)
	)

	publicKey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	if len(publicKey) != KeySize {
		return nil, fmt.Errorf("public key invalid length=%d", len(publicKey))
	}

	privateKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	if len(privateKey) != KeySize {
		return nil, fmt.Errorf("private key invalid length=%d", len(privateKey))
	}

	n := &Box{}
	copy(n.privateKey[:], privateKey)
	copy(n.publicKey[:], publicKey)

	return n, nil
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
