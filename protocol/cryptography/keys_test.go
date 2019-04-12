package cryptography

import (
	"testing"
)

func TestRandomKey(t *testing.T) {
	if RandomKey32().Equal(Key32{}) {
		t.Errorf("RandomKey32 returned empty key")
	}
}

func TestRandomKeyEmpty(t *testing.T) {
	if RandomKey32().Equal(RandomKey32()) {
		t.Errorf("RandomKey32 returned two the same keys")
	}
}

func TestKeyEqual(t *testing.T) {
	var (
		keyA = Key32{1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2,
		}
		keyB = Key32{1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2,
		}
	)

	if !keyA.Equal(keyB) {
		t.Errorf("Equal: keys are different %q %q, expected to be equal", keyA, keyB)
	}
}

func TestKeyNotEqual(t *testing.T) {
	var (
		keyA = Key32{1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2,
		}
		keyB = Key32{1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1,
		}
	)

	if keyA.Equal(keyB) {
		t.Errorf("Equal: keys are the same %q %q, expected to be different", keyA, keyB)
	}
}

func TestKeyString(t *testing.T) {
	var (
		keyA = Key32{1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2,
		}

		keyStr = "0102030405060708090001020304050607080900010203040506070809000102"
	)

	if keyA.String() != keyStr {
		t.Errorf("key String, got %q, expected %q", keyA.String(), keyStr)
	}
}

func TestKeyStringLength(t *testing.T) {
	var (
		keyA = Key32{1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2,
		}
	)

	if str := keyA.String(); len(str) != 64 {
		t.Errorf("key String length got %d, expected %d", len(str), 64)
	}
}

func TestKeyFromByteTooShort(t *testing.T) {
	byteKey := "121212"

	if _, err := Key32FromByte([]byte(byteKey)); err == nil {
		t.Errorf("Key32FromByte: expected error if key is too short")
	}
}

func TestKeyFromByte(t *testing.T) {
	var (
		byteKey = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2,
		}
	)

	key, err := Key32FromByte(byteKey)
	if err != nil {
		t.Errorf("Key32FromByte: unexpected error: %s", err)
	}

	if key.Equal(Key32{}) {
		t.Errorf("Key32FromByte returned empty key")
	}
}

func TestKeyFromByteNil(t *testing.T) {
	if _, err := Key32FromByte(nil); err == nil {
		t.Errorf("Key32FromByte: expected error if key is nil")
	}
}

func TestKeyFromStringEmpty(t *testing.T) {
	if _, err := Key32FromString(""); err == nil {
		t.Errorf("KeyFromString: expected error if key is empty")
	}
}

func TestKeyFromStringTooShort(t *testing.T) {
	tooShortKey := "121212"

	if _, err := Key32FromString(tooShortKey); err == nil {
		t.Errorf("KeyFromString(%q): expected error if key is too short", tooShortKey)
	}
}

func TestKeyFromStringInvalid(t *testing.T) {
	invalidKey := " oidad oai aasoif"

	if _, err := Key32FromString(invalidKey); err == nil {
		t.Errorf("KeyFromString(%q): expected error if key is invalid", invalidKey)
	}
}

func TestKeyFromString(t *testing.T) {
	var (
		expectedKey = Key32{1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2,
		}

		keyStr = "0102030405060708090001020304050607080900010203040506070809000102"
	)

	key, err := Key32FromString(keyStr)
	if err != nil {
		t.Errorf("KeyFromString(%q) unexpected error: %s", keyStr, err)
	}

	if !key.Equal(expectedKey) {
		t.Errorf("KeyFromString: expected %q, got %q", expectedKey, key)
	}
}
