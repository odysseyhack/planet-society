package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"testing"

	"github.com/andrewromanenco/gcrypt"
)

const (
	messageSize = 4096
)

func BenchmarkPrecompute(b *testing.B) {
	message := make([]byte, messageSize)
	if _, err := rand.Read(message); err != nil {
		b.Fatalf("read: failed: %s", err)
	}

	naclAlice, err := NewOneShotBox()
	if err != nil {
		b.Fatalf("NewOneShotBox failed: %s", err)
	}

	naclBob, err := NewOneShotBox()
	if err != nil {
		b.Fatalf("NewOneShotBox failed: %s", err)
	}

	sharedKeyAlice := naclAlice.Precompute(naclBob.GetPublicKey())
	sharedKeyBob := naclBob.Precompute(naclAlice.GetPublicKey())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encrypted, err := naclAlice.EncryptAfterPrecomputation(message, &sharedKeyAlice)
		if err != nil {
			b.Fatalf("encrypt failed: %s", err)
		}

		_, err = naclBob.DecryptAfterPrecomputation(encrypted, &sharedKeyBob)
		if err != nil {
			b.Fatalf("encrypt failed: %s", err)
		}
	}
}

func BenchmarkBox(b *testing.B) {
	message := make([]byte, messageSize)
	if _, err := rand.Read(message); err != nil {
		b.Fatalf("read: failed: %s", err)
	}

	naclAlice, err := NewOneShotBox()
	if err != nil {
		b.Errorf("Failed to prepare box one for test: %s", err)
	}

	naclBob, err := NewOneShotBox()
	if err != nil {
		b.Errorf("Failed to prepare box two for test: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encrypted, err := naclAlice.Encrypt([]byte(message), &naclBob.publicKey)
		if err != nil {
			b.Fatalf("failed to encrypt message: %s", err)
		}
		_, err = naclBob.Decrypt(encrypted, &naclAlice.publicKey)
		if err != nil {
			b.Fatalf("failed to decrypt message: %s", err)
		}
	}
}

func BenchmarkSecretBox(b *testing.B) {
	var (
		key       = RandomKey32()
		secretBox = NewSecretBox(key)
		message   = make([]byte, messageSize)
	)

	if _, err := rand.Read(message); err != nil {
		b.Fatalf("read: failed: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encrypted, err := secretBox.Encrypt([]byte(message))
		if err != nil {
			b.Fatalf("secretbox.Encrypt failed: %s", err)
		}

		_, err = secretBox.Decrypt(encrypted)
		if err != nil {
			b.Fatalf("secretbox.Decrypt failed: %s", err)
		}
	}
}

func BenchmarkNoHmacAES(b *testing.B) {
	var (
		key     = RandomKey32()
		message = make([]byte, messageSize)
	)

	block, err := aes.NewCipher(key[:])
	if err != nil {
		b.Fatalf("aes.NewCipher failed: %s", err)
	}

	if _, err := rand.Read(message); err != nil {
		b.Fatalf("read: failed: %s", err)
	}

	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		b.Fatalf("read: failed: %s", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		b.Fatalf("NewGCM failed: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := rand.Read(nonce); err != nil {
			b.Fatalf("read: failed: %s", err)
		}

		cipherText := aesgcm.Seal(nil, nonce, message, nil)
		_, err := aesgcm.Open(nil, nonce, cipherText, nil)
		if err != nil {
			b.Fatalf("Open failed: %s", err)
		}
	}
}

func BenchmarkHmacAES(b *testing.B) {
	var (
		key     = RandomKey32()
		message = make([]byte, messageSize)
	)

	if _, err := rand.Read(message); err != nil {
		b.Fatalf("read: failed: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ct, err := gcrypt.Encrypt(key[:], message)
		if err != nil {
			b.Fatalf("Encrypt failed: %s", err)
		}
		_, err = gcrypt.Decrypt(key[:], ct)
		if err != nil {
			b.Fatalf("Decrypt failed: %s", err)
		}
	}
}
