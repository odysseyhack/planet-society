package models

import (
	"fmt"
	"io"
	"strconv"

	"github.com/planet-platform/transaction-protocol-poc/cryptography"
)

type Key32 struct {
	Key cryptography.Key32
}

func (k Key32) MarshalGQL(w io.Writer) {
	data := strconv.Quote(k.Key.String())
	_, _ = w.Write([]byte(data))
}

func (k *Key32) UnmarshalGQL(v interface{}) error {
	switch v := v.(type) {
	case string:
		key, err := cryptography.Key32FromString(v)
		if err != nil {
			return err
		}
		k.Key = key
		return nil
	default:
		return fmt.Errorf("%T is invalid", v)
	}
}
