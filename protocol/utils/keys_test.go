package utils

import (
	"testing"

	"github.com/odysseyhack/planet-society/protocol/cryptography"
)

func TestKeys(t *testing.T) {
	exampleKey := cryptography.RandomKey32()
	if err := WriteKeyToDir(exampleKey); err != nil {
		t.Fatalf("failed to write key: %s", err)
	}

	gotKey, err := ReadKeyFromDir()
	if err != nil {
		t.Fatalf("failed to read key: %s", err)
	}

	if !gotKey.Equal(exampleKey) {
		t.Fatalf("key mismatch: %s %s", gotKey.String(), exampleKey.String())
	}
}
