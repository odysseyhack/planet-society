package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/odysseyhack/planet-society/protocol/cryptography"
)

const (
	keyFileName = "hackaton_key32.json"
	dir         = "/tmp/"
)

func WriteKeyToDir(key32 cryptography.Key32) error {
	data, err := json.Marshal(key32)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(dir, keyFileName), data, 0600)
}

func ReadKeyFromDir() (key32 cryptography.Key32, err error) {
	data, err := ioutil.ReadFile(filepath.Join(dir, keyFileName))
	if err != nil {
		return key32, err
	}

	err = json.Unmarshal(data, &key32)
	return key32, err
}

func CleanKey() error {
	return os.Remove(filepath.Join(dir, keyFileName))
}
