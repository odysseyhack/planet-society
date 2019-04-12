package cryptography

import "testing"

func TestGenerateSigner(t *testing.T) {
	signer, err := GenerateSigner()
	if err != nil {
		t.Fatalf("GenerateSigner failed: %s", err)
	}

	if signer == nil {
		t.Fatalf("GenerateSigner returned nil signer")
	}

}
