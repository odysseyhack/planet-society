package database

import (
	"github.com/boltdb/bolt"
	"github.com/odysseyhack/planet-society/protocol/cryptography"
	"github.com/odysseyhack/planet-society/protocol/models"
)

// bucketInitialize initializes root structure of the database
func (d *Database) bucketInitialize(tx *bolt.Tx) error {
	bucket, err := tx.CreateBucketIfNotExists([]byte(bucketPermissionsGranted))
	if err != nil {
		return err
	}

	if _, err := tx.CreateBucketIfNotExists([]byte(bucketIdentities)); err != nil {
		return err
	}

	bucket, err = tx.CreateBucketIfNotExists([]byte(personalDetailsBucket))
	if err != nil {
		return err
	}

	key := cryptography.RandomKey32()
	personalDetails := models.PersonalDetails{
		PublicKey: models.Key32{Key: d.keychain.MainPublicKey},
		ID:        key.String(),
	}

	return d.put(bucket, []byte(personalDetailsKey), &personalDetails)
}
