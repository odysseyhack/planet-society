/*
Package gcrypt is an util to work with aes 256 encryption. It supports key
derivation from user a provided password. The derivation either generates a new
salt or uses a provided one.

Salt must be unique for each password.

Encrypt uses 256 bit key (probably derivated from a user provided password) and
uses aes 256 with CFB mode. After a data is encrypted, HMAC is calculated and
appended to encrypted data.

Decrypt checks if HMAC is valid; return error if it is not.
*/
package gcrypt

import (
	"errors"
	"io"

	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/scrypt"
)

const defaultKeyLen = 32

// DerivateKey256 creates 256 bit key based on a password. Random salt is
// returned with the key.
func DerivateKey256(password string) ([]byte, []byte, error) {
	salt, err := generateSalt(defaultKeyLen)
	if err != nil {
		return nil, nil, err
	}
	key, err := DerivateKey256WithSalt(password, salt)
	if err != nil {
		return nil, nil, err
	}
	return key, salt, nil
}

// DerivateKey256WithSalt creates 256 bit key from provided password and salt.
func DerivateKey256WithSalt(password string, salt []byte) ([]byte, error) {
	if password == "" {
		return nil, errors.New("Empty pass")
	}
	if salt == nil || len(salt) != 32 {
		return nil, errors.New("Salt is not 256 bit")
	}
	// Recommended settings for scrypt
	key, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, defaultKeyLen)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func generateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// appendHMAC appends 32 bytes to data. Returns nil if no data is provided.
func appendHMAC(key, data []byte) []byte {
	if len(data) == 0 {
		return nil
	}
	macProducer := hmac.New(sha256.New, key)
	macProducer.Write(data)
	mac := macProducer.Sum(nil)
	return append(data, mac...)
}

// validateHMAC checks mac, and returns original data without mac bytes.
// Returns nil, if mac is not valid.
func validateHMAC(key, data []byte) []byte {
	if len(data) <= 32 {
		return nil
	}
	message := data[:len(data)-32]
	mac := data[len(data)-32:]
	macProducer := hmac.New(sha256.New, key)
	macProducer.Write(message)
	calculatedMac := macProducer.Sum(nil)
	if calculatedMac == nil {
		return nil
	}
	for i := 0; i < 32; i++ {
		if mac[i] != calculatedMac[i] {
			return nil
		}
	}
	return message
}

// Encrypt encrypts data with aes 256 and adds HMAC(EnM). Fails if key is not
// 256 bit or it data is empty.
func Encrypt(key, data []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("This is not a 256 bit key")
	}
	if len(data) == 0 {
		return nil, errors.New("Data is empty")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], data)
	return appendHMAC(key, ciphertext), nil
}

// Decrypt validates mac and returns decoded data.
func Decrypt(key, data []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("This is not a 256 bit key")
	}
	if len(data) == 0 {
		return nil, errors.New("Data is empty")
	}
	ciphertext := validateHMAC(key, data)
	if ciphertext == nil {
		return nil, errors.New("Invalid HMAC")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("Ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	result := ciphertext[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(result, result)
	return result, nil
}
