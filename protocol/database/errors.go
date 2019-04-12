package database

import "fmt"

// ErrKeyNotFound is returned when item with given key does not exist
func ErrKeyNotFound(key []byte) error {
	return fmt.Errorf("db: key %q not found", string(key))
}

// ErrBucketNotFound is returned when needed bucket does not exist
func ErrBucketNotFound(name string) error {
	return fmt.Errorf("db: bucket %q not found", string(name))
}

// ErrAlreadyExist is returned when item inside database already exist
func ErrAlreadyExist(name string) error {
	return fmt.Errorf("db: item %q already exist", string(name))
}
