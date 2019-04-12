package cryptography

import (
	"bytes"
	"testing"
)

func TestHash16Empty(t *testing.T) {
	var (
		hash      = Hash16([]byte{})
		emptyHash [16]byte
	)

	if bytes.Equal(emptyHash[:], hash[:]) {
		t.Errorf("Hash16 returned empty hash")
	}
}

func TestHash16NonEmpty(t *testing.T) {
	var (
		hash      = Hash16([]byte("some random text"))
		emptyHash [16]byte
	)

	if bytes.Equal(emptyHash[:], hash[:]) {
		t.Errorf("Hash16 returned empty hash")
	}
}
