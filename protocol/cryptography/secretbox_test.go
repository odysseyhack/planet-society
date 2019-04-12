package cryptography

import "testing"

func TestNewSecretBox(t *testing.T) {
	if NewSecretBox(RandomKey32()) == nil {
		t.Errorf("NewSecretBox returned nil handler")
	}
}

func TestSecretBox(t *testing.T) {
	var (
		key       = RandomKey32()
		secretBox = NewSecretBox(key)
		message   = "super secret message"
	)

	encrypted, err := secretBox.Encrypt([]byte(message))
	if err != nil {
		t.Fatalf("secretbox.Encrypt failed: %s", err)
	}

	if encrypted == nil {
		t.Fatalf("secretbox.Encrypt returned nil encrypted message")
	}

	decrypted, err := secretBox.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("secretbox.Decrypt failed: %s", err)
	}

	if decrypted == nil {
		t.Fatalf("secretbox.Decrypt returned nil decrypted message")
	}

	if string(decrypted) != message {
		t.Errorf("secretbox.Decrypt returned %q, expected %q", string(decrypted), message)
	}
}
