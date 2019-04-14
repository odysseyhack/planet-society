package cryptography

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestNewOneShotNaCl(t *testing.T) {
	var emptyKey Key32

	nacl, err := NewOneShotBox()

	if err != nil {
		t.Fatalf("NewOneShotBox failed: %s", err)
	}

	if nacl == nil {
		t.Fatalf("NewOneShotBox returned nil NaCl")
	}

	if bytes.Equal(emptyKey[:], nacl.privateKey[:]) {
		t.Fatalf("NewOneShotBox returned NaCl with empty private key")
	}

	if bytes.Equal(emptyKey[:], nacl.publicKey[:]) {
		t.Fatalf("NewOneShotBox returned NaCl with empty public key")
	}
}

func TestManySizeEncrypt(t *testing.T) {
	naclAlice, err := NewOneShotBox()
	if err != nil {
		t.Fatalf("NewOneShotBox failed: %s", err)
	}

	naclBob, err := NewOneShotBox()
	if err != nil {
		t.Fatalf("NewOneShotBox failed: %s", err)
	}

	for i := uint(0); i <= 14; i++ {
		dataSize := (1 << i) & (^uint(1))
		randMsg := make([]byte, dataSize)
		rand.Read(randMsg)

		encryptedData, err := naclAlice.Encrypt(randMsg, &naclBob.publicKey)
		if err != nil {
			t.Fatalf("Encrypt failed: %s", err)
		}

		decrypted, err := naclBob.Decrypt(encryptedData, &naclAlice.publicKey)
		if err != nil {
			t.Fatalf("Decrypt failed: %s", err)
		}

		if !bytes.Equal(decrypted, randMsg) {
			t.Fatalf("at size=%d , mismatch between data and decrypted data", dataSize)
		}

	}
}

func TestNewBox(t *testing.T) {
	var (
		privateKey = RandomKey32()
		publicKey  = RandomKey32()
	)

	box := NewBox(privateKey, publicKey)

	if box == nil {
		t.Fatalf("NewBox returned nil")
	}

	if !box.publicKey.Equal(publicKey) {
		t.Errorf("NewBox set wrong public key")
	}

	if !box.privateKey.Equal(privateKey) {
		t.Errorf("NewBox set wrong private key")
	}
}

func TestBoxEncryptDecryptSelf(t *testing.T) {
	var (
		message = []byte("Some message for encryption")
	)

	box, err := NewOneShotBox()
	if err != nil {
		t.Errorf("Failed to prepare box for test: %s", err)
	}

	encrypted, err := BoxEncrypt(message, &box.publicKey, &box.privateKey)
	if err != nil {
		t.Errorf("BoxEncrypt failed: %s", err)
	}

	if encrypted == nil {
		t.Errorf("BoxEncrypt returned nil encryption message")
	}

	if bytes.Equal(message, encrypted) {
		t.Errorf("message is not encrypted")
	}

	decrypted, err := BoxDecrypt(encrypted, &box.privateKey, &box.publicKey)
	if err != nil {
		t.Errorf("BoxEncrypt failed: %s", err)
	}

	if decrypted == nil {
		t.Errorf("BoxDecrypt returned nil decrypted message")
	}

	if !bytes.Equal(message, decrypted) {
		t.Errorf("decrypted message is not the same as original")
	}
}

func TestBoxEncryptDecrypt(t *testing.T) {
	var (
		message = []byte("Some message for encryption")
	)

	box1, err := NewOneShotBox()
	if err != nil {
		t.Errorf("Failed to prepare box one for test: %s", err)
	}

	box2, err := NewOneShotBox()
	if err != nil {
		t.Errorf("Failed to prepare box two for test: %s", err)
	}

	encrypted, err := BoxEncrypt(message, &box2.publicKey, &box1.privateKey)
	if err != nil {
		t.Errorf("BoxEncrypt failed: %s", err)
	}

	if encrypted == nil {
		t.Errorf("BoxEncrypt returned nil encryption message")
	}

	if bytes.Equal(message, encrypted) {
		t.Errorf("message is not encrypted")
	}

	decrypted, err := BoxDecrypt(encrypted, &box2.privateKey, &box1.publicKey)
	if err != nil {
		t.Errorf("BoxEncrypt failed: %s", err)
	}

	if decrypted == nil {
		t.Errorf("BoxDecrypt returned nil decrypted message")
	}

	if !bytes.Equal(message, decrypted) {
		t.Errorf("decrypted message is not the same as original")
	}
}

func TestPrecomputation(t *testing.T) {
	var (
		msg = []byte("secret message")
	)
	naclAlice, err := NewOneShotBox()
	if err != nil {
		t.Fatalf("NewOneShotBox failed: %s", err)
	}

	naclBob, err := NewOneShotBox()
	if err != nil {
		t.Fatalf("NewOneShotBox failed: %s", err)
	}

	sharedKeyAlice := naclAlice.Precompute(naclBob.GetPublicKey())
	sharedKeyBob := naclBob.Precompute(naclAlice.GetPublicKey())

	encrypted, err := naclAlice.EncryptAfterPrecomputation(msg, &sharedKeyAlice)
	if err != nil {
		t.Fatalf("encrypt failed: %s", err)
	}

	if encrypted == nil {
		t.Fatalf("encrypt returned nil encrypted message")
	}

	decrypted, err := naclBob.DecryptAfterPrecomputation(encrypted, &sharedKeyBob)
	if err != nil {
		t.Fatalf("encrypt failed: %s", err)
	}

	if !bytes.Equal(msg, decrypted) {
		t.Errorf("encrypted message is not equal original one")
	}
}
