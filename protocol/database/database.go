package database

import (
	"bytes"
	"encoding/gob"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/odysseyhack/planet-society/protocol/cryptography"
	"github.com/odysseyhack/planet-society/protocol/models"
)

// Db root tree:
//   -> permissions [Key32=transaction ID]
//       -> permission 1
//       -> permission 2
//   -> identities [identity_ID]
//          -> contacts [ID]
//          -> addresses [ID]
//              -> payment_cards [ID]
//              -> passports [ID]
//              -> identity_document [ID]

type Database struct {
	db       *bolt.DB
	keychain *cryptography.Keychain
}

// LoadDatabase loads database from the file and returns handler to it.
// If the database is not existing it will be created and initialized with buckets.
func LoadDatabase(filePath string, keychain *cryptography.Keychain) (*Database, error) {
	directory, _ := filepath.Split(filePath)
	if err := os.MkdirAll(directory, 0700); err != nil {
		return nil, err
	}

	_, statErr := os.Stat(filePath)

	db, err := bolt.Open(filePath, 0600, nil)
	if err != nil {
		return nil, err
	}

	database := &Database{
		db:       db,
		keychain: keychain,
	}
	if err := database.initialize(os.IsNotExist(statErr)); err != nil {
		return nil, err
	}

	return database, nil
}

// Initialize initializes database with proper buckets
func (d *Database) initialize(run bool) error {
	if run {
		return d.db.Update(d.bucketInitialize)
	}
	return nil
}

// Close closes database
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// newID generates unique 32 bytes long identifier for item inside database
func (d *Database) newID() string {
	key := cryptography.RandomKey32()
	return key.String()
}

// put serializes and inserts object inside database
func (d *Database) put(bucket *bolt.Bucket, key []byte, object interface{}) error {
	var buffer bytes.Buffer

	if err := gob.NewEncoder(&buffer).Encode(object); err != nil {
		return err
	}

	return bucket.Put(key, buffer.Bytes())
}

// get reads and deserialize object from the database
func (d *Database) get(bucket *bolt.Bucket, key []byte, object interface{}) error {
	raw := bucket.Get(key)
	if raw == nil {
		return ErrKeyNotFound(key)
	}

	return d.decode(raw, object)
}

// decode deserialize object from the database
func (d *Database) decode(value []byte, object interface{}) error {
	if err := gob.NewDecoder(bytes.NewBuffer(value)).Decode(object); err != nil {
		return err
	}
	return nil
}

func updatePersonalDetails(personal *models.PersonalDetailsInput, updatedDetail *models.PersonalDetails) {
	if personal.Name != nil {
		updatedDetail.Name = *personal.Name
	}

	if personal.Surname != nil {
		updatedDetail.Surname = *personal.Surname
	}

	if personal.Country != nil {
		updatedDetail.Country = *personal.Country
	}

	if personal.BirthDate != nil {
		updatedDetail.BirthDate = *personal.BirthDate
	}
}

// PersonalDetailsAdd adds personal details to the database
func (d *Database) PersonalDetailsUpdate(personal models.PersonalDetailsInput) (updatedDetail models.PersonalDetails, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(personalDetailsKey))
		if bucket == nil {
			return ErrKeyNotFound([]byte(personalDetailsKey))
		}

		if err := d.get(bucket, []byte(personalDetailsKey), &updatedDetail); err != nil {
			return err
		}

		updatePersonalDetails(&personal, &updatedDetail)

		return d.put(bucket, []byte(personalDetailsKey), &updatedDetail)
	})

	return updatedDetail, err
}

// PersonalDetails returns personal details
func (d *Database) PersonalDetails() (details models.PersonalDetails, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(personalDetailsKey))
		if bucket == nil {
			return ErrKeyNotFound([]byte(personalDetailsKey))
		}
		return d.get(bucket, []byte(personalDetailsKey), &details)
	})

	details.Email = "example@example.com"
	details.Bsn = "128937192837"

	return details, err
}

func (d *Database) createNewIdentitySubBuckets(identitiesBucket *bolt.Bucket, newIdentity *models.Identity) error {
	newIdentityBucket, err := identitiesBucket.CreateBucketIfNotExists([]byte(newIdentity.ID))
	if err != nil {
		return err
	}

	if err := d.put(newIdentityBucket, []byte(identityMetadataKey), &newIdentity); err != nil {
		return err
	}

	if _, err := newIdentityBucket.CreateBucketIfNotExists([]byte(bucketContacts)); err != nil {
		return err
	}

	if _, err := newIdentityBucket.CreateBucketIfNotExists([]byte(bucketAddress)); err != nil {
		return err
	}

	if _, err := newIdentityBucket.CreateBucketIfNotExists([]byte(bucketPassports)); err != nil {
		return err
	}

	if _, err := newIdentityBucket.CreateBucketIfNotExists([]byte(bucketPaymentCards)); err != nil {
		return err
	}

	if _, err := newIdentityBucket.CreateBucketIfNotExists([]byte(bucketIdentityDocuments)); err != nil {
		return err
	}
	return nil
}

// IdentityAdd adds new identities
func (d *Database) IdentityAdd(identity models.IdentityInput) (newIdentity models.Identity, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		identitiesBucket := tx.Bucket([]byte(bucketIdentities))
		if identitiesBucket == nil {
			return ErrBucketNotFound(bucketIdentities)
		}

		newIdentity = models.Identity{
			ID:          d.newID(),
			DisplayName: identity.DisplayName,
		}

		if err := d.createNewIdentitySubBuckets(identitiesBucket, &newIdentity); err != nil {
			return err
		}

		return nil
	})

	return newIdentity, err
}

// IdentityList returns list of stored identities
func (d *Database) IdentityList() (list []models.Identity, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		identitiesBucket := tx.Bucket([]byte(bucketIdentities))
		if identitiesBucket == nil {
			return ErrBucketNotFound(bucketIdentities)
		}

		return identitiesBucket.ForEach(func(k, v []byte) error {
			bucket := identitiesBucket.Bucket(k)
			if bucket == nil {
				return ErrBucketNotFound(string(k))
			}

			var metadata models.Identity
			if err := d.get(bucket, []byte(identityMetadataKey), &metadata); err != nil {
				return err
			}
			list = append(list, metadata)
			return nil
		})
	})

	return list, err
}

func (d *Database) contactInputToContact(contact *models.ContactInput, newContact *models.Contact) {
	newContact.ID = d.newID()
	newContact.DisplayName = contact.DisplayName
	newContact.PublicKey = contact.PublicKey
	newContact.SignatureKey = contact.SignatureKey

	if contact.Name != nil {
		newContact.Name = *contact.Name
	}

	if contact.Surname != nil {
		newContact.Surname = *contact.Surname
	}

	if contact.Country != nil {
		newContact.Country = *contact.Country
	}
}

// ContactAdd adds new contact to database
func (d *Database) ContactAdd(contact models.ContactInput) (newContact models.Contact, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		identitiesBucket := tx.Bucket([]byte(bucketIdentities))
		if identitiesBucket == nil {
			return ErrBucketNotFound(bucketIdentities)
		}

		idBucket := identitiesBucket.Bucket([]byte(contact.Identity))
		if idBucket == nil {
			return ErrBucketNotFound(contact.Identity)
		}
		contactBucket := idBucket.Bucket([]byte(bucketContacts))
		if contactBucket == nil {
			return ErrBucketNotFound(bucketContacts)
		}

		d.contactInputToContact(&contact, &newContact)
		return d.put(contactBucket, []byte(newContact.ID), &newContact)
	})

	return newContact, err
}

func (d *Database) collectContacts(list *[]models.Contact, bucket *bolt.Bucket) error {
	return bucket.ForEach(func(k, v []byte) error {
		var contact models.Contact
		if err := d.decode(v, &contact); err != nil {
			return err
		}
		*list = append(*list, contact)
		return nil
	})
}

// ContactList lists contacts in given identity
func (d *Database) ContactList(identity string) (list []models.Contact, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		identitiesBucket := tx.Bucket([]byte(bucketIdentities))
		if identitiesBucket == nil {
			return ErrBucketNotFound(bucketIdentities)
		}

		identityBucket := identitiesBucket.Bucket([]byte(identity))
		if identityBucket == nil {
			return ErrBucketNotFound(identity)
		}

		contactBucket := identityBucket.Bucket([]byte(bucketContacts))
		if contactBucket == nil {
			return ErrBucketNotFound(bucketContacts)
		}
		return d.collectContacts(&list, contactBucket)

	})
	return list, err
}

func (d *Database) addressInputToAddress(inputAddress *models.AddressInput, address *models.Address) {
	address.ID = d.newID()
	address.DisplayName = inputAddress.DisplayName

	if inputAddress.Country != nil {
		address.Country = *inputAddress.Country
	}
	if inputAddress.Street != nil {
		address.Street = *inputAddress.Street
	}
	if inputAddress.City != nil {
		address.City = *inputAddress.City
	}
}

// AddressAdd adds new address
func (d *Database) AddressAdd(addresses models.AddressInput) (added models.Address, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		identitiesBucket := tx.Bucket([]byte(bucketIdentities))
		if identitiesBucket == nil {
			return ErrBucketNotFound(bucketIdentities)
		}
		idBucket := identitiesBucket.Bucket([]byte(addresses.Identity))
		if idBucket == nil {
			return ErrBucketNotFound(addresses.Identity)
		}
		addressBucket := idBucket.Bucket([]byte(bucketAddress))
		if addressBucket == nil {
			return ErrBucketNotFound(bucketAddress)
		}
		d.addressInputToAddress(&addresses, &added)
		return d.put(addressBucket, []byte(added.ID), &added)
	})
	return added, err
}

func (d *Database) collectAddresses(list *[]models.Address, bucket *bolt.Bucket) error {
	return bucket.ForEach(func(k, v []byte) error {
		var address models.Address
		if err := d.decode(v, &address); err != nil {
			return err
		}
		*list = append(*list, address)
		return nil
	})
}

// AddressList lists known addresses
func (d *Database) AddressList(identity string) (list []models.Address, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		identitiesBucket := tx.Bucket([]byte(bucketIdentities))
		if identitiesBucket == nil {
			return ErrBucketNotFound(bucketIdentities)
		}

		identityBucket := identitiesBucket.Bucket([]byte(identity))
		if identityBucket == nil {
			return ErrBucketNotFound(identity)
		}

		addrBucket := identityBucket.Bucket([]byte(bucketAddress))
		if addrBucket == nil {
			return ErrBucketNotFound(bucketAddress)
		}
		return d.collectAddresses(&list, addrBucket)

	})
	return list, err
}

func (d *Database) paymentCardInputToPaymentCard(paymentCard *models.PaymentCardInput, added *models.PaymentCard) {
	added.ID = d.newID()
	added.Expiration = paymentCard.Expiration
	added.DisplayName = paymentCard.DisplayName
	added.Currency = paymentCard.Currency
	added.SecurityCode = paymentCard.SecurityCode
	added.Number = paymentCard.Number
	added.Surname = paymentCard.Surname
	added.Name = paymentCard.Name
}

// PaymentCardAdd adds new payment card
func (d *Database) PaymentCardAdd(paymentCard models.PaymentCardInput) (added models.PaymentCard, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		identitiesBucket := tx.Bucket([]byte(bucketIdentities))
		if identitiesBucket == nil {
			return ErrBucketNotFound(bucketIdentities)
		}
		// check if identities exist
		idBucket := identitiesBucket.Bucket([]byte(paymentCard.Identity))
		if idBucket == nil {
			return ErrBucketNotFound(paymentCard.Identity)
		}
		paymentCardBucket := idBucket.Bucket([]byte(bucketPaymentCards))
		if paymentCardBucket == nil {
			return ErrBucketNotFound(bucketPaymentCards)
		}

		d.paymentCardInputToPaymentCard(&paymentCard, &added)
		return d.put(paymentCardBucket, []byte(added.ID), &added)
	})
	return added, err
}

func (d *Database) collectPaymentCards(list *[]models.PaymentCard, contactBucket *bolt.Bucket) error {
	return contactBucket.ForEach(func(k, v []byte) error {
		var paymentCard models.PaymentCard
		if err := d.decode(v, &paymentCard); err != nil {
			return err
		}
		*list = append(*list, paymentCard)
		return nil
	})
}

// PaymentCardList lists known payment cards for given identity
func (d *Database) PaymentCardList(identity string) (list []models.PaymentCard, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		identitiesBucket := tx.Bucket([]byte(bucketIdentities))
		if identitiesBucket == nil {
			return ErrBucketNotFound(bucketIdentities)
		}
		identityBucket := identitiesBucket.Bucket([]byte(identity))
		if identityBucket == nil {
			return ErrBucketNotFound(identity)
		}
		contactBucket := identityBucket.Bucket([]byte(bucketPaymentCards))
		if contactBucket == nil {
			return ErrBucketNotFound(bucketContacts)
		}
		return d.collectPaymentCards(&list, contactBucket)
	})
	return list, err
}

func (d *Database) passportInputToPassport(passport *models.PassportInput, added *models.Passport) {
	added.ID = d.newID()
	added.Expiration = passport.Expiration
	added.DisplayName = passport.DisplayName
	added.Number = passport.Number
	added.Surname = passport.Surname
	added.Name = passport.Name
	added.Country = passport.Country
}

// PassportAdd adds new passport
func (d *Database) PassportAdd(passport models.PassportInput) (added models.Passport, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		identitiesBucket := tx.Bucket([]byte(bucketIdentities))
		if identitiesBucket == nil {
			return ErrBucketNotFound(bucketIdentities)
		}

		idBucket := identitiesBucket.Bucket([]byte(passport.Identity))
		if idBucket == nil {
			return ErrBucketNotFound(passport.Identity)
		}

		passportBucket := idBucket.Bucket([]byte(bucketPassports))
		if passportBucket == nil {
			return ErrBucketNotFound(bucketPassports)
		}

		d.passportInputToPassport(&passport, &added)
		return d.put(passportBucket, []byte(added.ID), &added)
	})
	return added, err
}

func (d *Database) collectPassport(list *[]models.Passport, passportBucket *bolt.Bucket) error {
	return passportBucket.ForEach(func(k, v []byte) error {
		var passport models.Passport
		if err := d.decode(v, &passport); err != nil {
			return err
		}
		*list = append(*list, passport)
		return nil
	})
}

// PassportList lists known passports for given identity
func (d *Database) PassportList(identity string) (list []models.Passport, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		identitiesBucket := tx.Bucket([]byte(bucketIdentities))
		if identitiesBucket == nil {
			return ErrBucketNotFound(bucketIdentities)
		}

		identityBucket := identitiesBucket.Bucket([]byte(identity))
		if identityBucket == nil {
			return ErrBucketNotFound(identity)
		}

		passportBucket := identityBucket.Bucket([]byte(bucketPassports))
		if passportBucket == nil {
			return ErrBucketNotFound(bucketContacts)
		}
		return d.collectPassport(&list, passportBucket)
	})
	return list, err
}

func (d *Database) identityDocumentInputToDocument(identityDocument *models.IdentityDocumentInput, document *models.IdentityDocument) {
	document.ID = d.newID()
	document.Expiration = identityDocument.Expiration
	document.DisplayName = identityDocument.DisplayName
	document.Number = identityDocument.Number
	document.Surname = identityDocument.Surname
	document.Name = identityDocument.Name
	document.Country = identityDocument.Country
}

// IdentityDocumentAdd adds new identity document
func (d *Database) IdentityDocumentAdd(identityDocument models.IdentityDocumentInput) (document models.IdentityDocument, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		identitiesBucket := tx.Bucket([]byte(bucketIdentities))
		if identitiesBucket == nil {
			return ErrBucketNotFound(bucketIdentities)
		}

		idBucket := identitiesBucket.Bucket([]byte(identityDocument.Identity))
		if idBucket == nil {
			return ErrBucketNotFound(identityDocument.Identity)
		}

		identityDocumentBucket := idBucket.Bucket([]byte(bucketIdentityDocuments))
		if identityDocumentBucket == nil {
			return ErrBucketNotFound(bucketIdentityDocuments)
		}

		d.identityDocumentInputToDocument(&identityDocument, &document)
		return d.put(identityDocumentBucket, []byte(document.ID), &document)
	})

	return document, err
}

func (d *Database) collectIdentityDocument(list *[]models.IdentityDocument, idsBucket *bolt.Bucket) error {
	return idsBucket.ForEach(func(k, v []byte) error {
		var identityDocument models.IdentityDocument
		if err := d.decode(v, &identityDocument); err != nil {
			return err
		}
		*list = append(*list, identityDocument)
		return nil
	})
}

// IdentityDocumentList lists known identity documents for given identity
func (d *Database) IdentityDocumentList(identity string) (list []models.IdentityDocument, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		identitiesBucket := tx.Bucket([]byte(bucketIdentities))
		if identitiesBucket == nil {
			return ErrBucketNotFound(bucketIdentities)
		}

		identityBucket := identitiesBucket.Bucket([]byte(identity))
		if identityBucket == nil {
			return ErrBucketNotFound(identity)
		}

		contactBucket := identityBucket.Bucket([]byte(bucketIdentityDocuments))
		if contactBucket == nil {
			return ErrBucketNotFound(bucketContacts)
		}
		return d.collectIdentityDocument(&list, contactBucket)
	})
	return list, err
}

func (d *Database) collectPermissions(list *[]models.Permission, bucket *bolt.Bucket) error {
	return bucket.ForEach(func(k, v []byte) error {
		var permission models.Permission
		if err := d.decode(v, &permission); err != nil {
			return err
		}
		*list = append(*list, permission)
		return nil
	})
}

func (d *Database) PermissionList() (list []models.Permission, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		permissionsBucket := tx.Bucket([]byte(bucketPermissionsGranted))
		if permissionsBucket == nil {
			return ErrBucketNotFound(bucketPermissionsGranted)
		}
		return d.collectPermissions(&list, permissionsBucket)

	})
	return list, err
}

// IdentityDel removed identity from database
// if identity with given id does not exist nil error is removed
func (d *Database) IdentityDel(id string) (removedID string, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketIdentities))
		if bucket == nil {
			return ErrBucketNotFound(bucketIdentities)
		}

		return bucket.Delete([]byte(id))
	})
	return id, err
}

func (d *Database) AddressDel(id string) (removedID string, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		idsForAddrDel := tx.Bucket([]byte(bucketIdentities))
		if idsForAddrDel == nil {
			return ErrBucketNotFound(bucketIdentities)
		}

		return idsForAddrDel.ForEach(func(k, v []byte) error {
			identityBucket := idsForAddrDel.Bucket(k)
			if identityBucket == nil {
				return nil
			}

			addressesBucket := identityBucket.Bucket([]byte(bucketAddress))
			if addressesBucket == nil {
				return nil
			}

			return addressesBucket.Delete([]byte(id))
		})
	})
	return id, err
}

func (d *Database) PassportDel(id string) (removedID string, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		idForPassportDel := tx.Bucket([]byte(bucketIdentities))
		if idForPassportDel == nil {
			return ErrBucketNotFound(bucketIdentities)
		}

		return idForPassportDel.ForEach(func(k, v []byte) error {
			identityBucket := idForPassportDel.Bucket(k)
			if identityBucket == nil {
				return nil
			}

			passportsBucket := identityBucket.Bucket([]byte(bucketPassports))
			if passportsBucket == nil {
				return nil
			}

			return passportsBucket.Delete([]byte(id))
		})
	})
	return id, err
}

func (d *Database) PaymentCardDel(id string) (removedID string, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		idsForPaymentDel := tx.Bucket([]byte(bucketIdentities))
		if idsForPaymentDel == nil {
			return ErrBucketNotFound(bucketIdentities)
		}

		return idsForPaymentDel.ForEach(func(k, v []byte) error {
			identityBucket := idsForPaymentDel.Bucket(k)
			if identityBucket == nil {
				return nil
			}

			paymentBucket := identityBucket.Bucket([]byte(bucketPaymentCards))
			if paymentBucket == nil {
				return nil
			}

			return paymentBucket.Delete([]byte(id))
		})
	})
	return id, err
}

func (d *Database) IdentityDocumentDel(id string) (removedID string, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		idForDocumentDel := tx.Bucket([]byte(bucketIdentities))
		if idForDocumentDel == nil {
			return ErrBucketNotFound(bucketIdentities)
		}

		return idForDocumentDel.ForEach(func(k, v []byte) error {
			identityBucket := idForDocumentDel.Bucket(k)
			if identityBucket == nil {
				return nil
			}

			documentBucket := identityBucket.Bucket([]byte(bucketIdentityDocuments))
			if identityBucket == nil {
				return nil
			}

			return documentBucket.Delete([]byte(id))
		})
	})
	return id, err
}

func (d *Database) ContactDel(id string) (removedID string, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		idForContactDel := tx.Bucket([]byte(bucketIdentities))
		if idForContactDel == nil {
			return ErrBucketNotFound(bucketIdentities)
		}

		return idForContactDel.ForEach(func(k, v []byte) error {
			identityBucket := idForContactDel.Bucket(k)
			if identityBucket == nil {
				return nil
			}

			contactsBucket := identityBucket.Bucket([]byte(bucketContacts))
			if contactsBucket == nil {
				return nil
			}

			return contactsBucket.Delete([]byte(id))
		})
	})
	return id, err
}

func (d *Database) permissionInputToPermission(permission, added *models.Permission) {
	added.TransactionID = permission.TransactionID
	added.Expiration = permission.Expiration
	added.Description = permission.Description
	added.Title = permission.Title
	added.RequesterPublicKey = permission.RequesterPublicKey
	added.RequesterSignatureKey = permission.RequesterSignatureKey
	added.RequesterSignature = permission.RequesterSignature
	added.ResponderSignature = permission.ResponderSignature
	added.PermissionNodes = permission.PermissionNodes
	added.Revokable = permission.Revokable
	added.ID = d.newID()
	added.LawApplying = permission.LawApplying
	added.LegalReliationships = permission.LegalReliationships
}

func (d *Database) PermissionAdd(permission models.Permission) (added models.Permission, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		permissionBucket := tx.Bucket([]byte(bucketPermissionsGranted))
		if permissionBucket == nil {
			return ErrBucketNotFound(bucketIdentities)
		}

		d.permissionInputToPermission(&permission, &added)

		if err := d.put(permissionBucket, []byte(added.ID), &added); err != nil {
			return err
		}

		return nil
	})
	return added, err
}
