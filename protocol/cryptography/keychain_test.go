package cryptography

import "testing"

func TestOneShotKeychain(t *testing.T) {
	keychain, err := OneShotKeychain()
	if err != nil {
		t.Errorf("OneShotKeychain failed: %s", err)
	}

	if keychain == nil {
		t.Errorf("OneShotKeychain returned nil keychain")
	}
}
