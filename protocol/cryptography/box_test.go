package cryptography

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestGenerateKeysSafe(t *testing.T) {
	var (
		path           = "/tmp/generate_key_test"
		privateKeyPath = fmt.Sprintf("%s/%s", path, privateKeyName)
		publicKeyPath  = fmt.Sprintf("%s/%s", path, publicKeyName)
	)

	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("failed to create dir for test: %s", err)
	}

	defer func() {
		if err := os.RemoveAll(path); err != nil {
			t.Errorf("failed to clean dir for test: %s", err)
		}
	}()

	if err := GenerateKeysToFile(path); err != nil {
		t.Errorf("GenerateKeysToFile failed: %s", err)
	}

	if _, err := os.Stat(privateKeyPath); err != nil {
		t.Fatalf("GenerateKeys: private key was not generated")
	}

	if _, err := os.Stat(publicKeyPath); err != nil {
		t.Fatalf("GenerateKeys: public key was not generated")
	}

	privateKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		t.Fatalf("failed to read private key: %s", err)
	}

	publicKey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		t.Fatalf("failed to read public key: %s", err)
	}

	if len(privateKey) != KeySize {
		t.Errorf("private key len=%d, expected %d", len(privateKey), KeySize)
	}

	if len(publicKey) != KeySize {
		t.Errorf("public key len=%d, expected %d", len(publicKey), KeySize)
	}
}

func TestGenerateKeysSafeError(t *testing.T) {
	var (
		path = "/tmp/generate_key_test"
	)

	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("failed to create dir for test: %s", err)
	}

	defer func() {
		if err := os.RemoveAll(path); err != nil {
			t.Errorf("failed to clean dir for test: %s", err)
		}
	}()

	if err := GenerateKeysToFile(path); err != nil {
		t.Errorf("GenerateKeysToFile failed: %s", err)
	}

	if err := GenerateKeysToFile(path); err == nil {
		t.Errorf("expected error when calling GenerateKeysToFile for the second time")
	}
}

func TestLoadGenerated(t *testing.T) {
	var (
		path     = "/tmp/load_key_test"
		emptyKey [KeySize]byte
	)

	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("failed to create dir for test: %s", err)
	}

	defer func() {
		if err := os.RemoveAll(path); err != nil {
			t.Errorf("failed to clean dir for test: %s", err)
		}
	}()

	if err := GenerateKeysToFile(path); err != nil {
		t.Errorf("GenerateKeysToFile failed: %s", err)
	}

	nacl, err := BoxLoad(path)
	if err != nil {
		t.Fatalf("Box.Load failed: %s", err)
	}

	if bytes.Equal(nacl.privateKey[:], emptyKey[:]) {
		t.Errorf("Box.Load: private key is empty")
	}

	if bytes.Equal(nacl.publicKey[:], emptyKey[:]) {
		t.Errorf("Box.Load: public key is empty")
	}
}

func TestEncryption(t *testing.T) {
	var (
		pathAlice = "/tmp/naclAlice_test"
		pathBob   = "/tmp/naclBob_test"
		message   = "This is message which should be encrypted"
	)

	if err := os.MkdirAll(pathAlice, 0755); err != nil {
		t.Fatalf("failed to create dir for test: %s", err)
	}

	if err := os.MkdirAll(pathBob, 0755); err != nil {
		t.Fatalf("failed to create dir for test: %s", err)
	}

	defer func() {
		if err := os.RemoveAll(pathAlice); err != nil {
			t.Errorf("failed to clean dir for test: %s", err)
		}
		if err := os.RemoveAll(pathBob); err != nil {
			t.Errorf("failed to clean dir for test: %s", err)
		}
	}()

	if err := GenerateKeysToFile(pathAlice); err != nil {
		t.Errorf("GenerateKeysToFile failed: %s", err)
	}

	if err := GenerateKeysToFile(pathBob); err != nil {
		t.Errorf("GenerateKeysToFile failed: %s", err)
	}

	naclAlice, err := BoxLoad(pathAlice)
	if err != nil {
		t.Fatalf("BoxLoad failed: %s", err)
	}

	naclBob, err := BoxLoad(pathBob)
	if err != nil {
		t.Fatalf("BoxLoad failed: %s", err)
	}

	encrypted, err := naclAlice.Encrypt([]byte(message), &naclBob.publicKey)
	if err != nil {
		t.Fatalf("failed to encrypt message: %s", err)
	}

	if len(encrypted) == 0 {
		t.Fatalf("encrypted message is empty")
	}

	decrypted, err := naclBob.Decrypt(encrypted, &naclAlice.publicKey)
	if err != nil {
		t.Fatalf("failed to decrypt message: %s", err)
	}

	if string(decrypted) != message {
		t.Errorf("decrypted message is not equal to original one")
	}
}

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
