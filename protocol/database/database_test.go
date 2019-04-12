package database

import (
	"os"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/odysseyhack/planet-society/protocol/cryptography"
	"github.com/odysseyhack/planet-society/protocol/models"
)

func TestLoadDatabaseExist(t *testing.T) {
	const (
		fileName = "/tmp/test_dir_i2i/database/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}

	}()

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Errorf("database file %q does not exist", fileName)
	}
}

func TestLoadDatabaseInitialized(t *testing.T) {
	const (
		fileName = "/tmp/test_dir_i2i/initialize/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	bucketExist(db.db, bucketPermissionsGranted, t)
	bucketExist(db.db, bucketIdentities, t)
	bucketExist(db.db, personalDetailsKey, t)
}

func TestLoadDatabaseInitializedTwice(t *testing.T) {
	const (
		fileName = "/tmp/test_dir_i2i/initialize/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	if err := db.Close(); err != nil {
		t.Errorf("Close failed: %s", err)
	}

	dbSecond, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	bucketExist(dbSecond.db, bucketPermissionsGranted, t)
	bucketExist(dbSecond.db, bucketIdentities, t)
	bucketExist(dbSecond.db, personalDetailsKey, t)

	if err := dbSecond.Close(); err != nil {
		t.Errorf("Close failed: %s", err)
	}
}

func TestNewID(t *testing.T) {
	var database Database
	id1 := database.newID()
	id2 := database.newID()

	if id1 == id2 {
		t.Errorf("newID: returned the same identifiers")
	}
}

func TestPersonalDetail(t *testing.T) {
	const (
		fileName = "/tmp/test_dir_i2i/pd1/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	details, err := db.PersonalDetails()
	if err != nil {
		t.Fatalf("PersonalDetails() failed: %s", err)
	}

	if !details.PublicKey.Key.Equal(wallet.MainPublicKey) {
		t.Errorf("PersonalDetail.PublicKey: got %q, expected %q",
			details.PublicKey.Key.String(),
			wallet.MainPublicKey.String())
	}
}

func TestPersonalDetailUpdateNoChange(t *testing.T) {
	const (
		fileName = "/tmp/test_dir_i2i/pd1/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	_, err = db.PersonalDetailsUpdate(models.PersonalDetailsInput{})
	if err != nil {
		t.Fatalf("PersonalDetailsUpdate() failed: %s", err)
	}
}

func TestPersonalDetailUpdateName(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd1/file.db"
		name     = "Tom"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	updateReturn, err := db.PersonalDetailsUpdate(models.PersonalDetailsInput{Name: &name})
	if err != nil {
		t.Fatalf("PersonalDetailsUpdate() failed: %s", err)
	}

	if updateReturn.Name != name {
		t.Errorf("PersonalDetailsUpdate expected to return details with proper name")
	}

	details, err := db.PersonalDetails()
	if err != nil {
		t.Fatalf("PersonalDetails() failed: %s", err)
	}

	if details.Name != name {
		t.Errorf("PersonalDetail.Name: got %q, expected %q", details.Name, name)
	}
}

func TestPersonalDetailUpdateSurname(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd1/file.db"
		surname  = "Hanks"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	updateReturn, err := db.PersonalDetailsUpdate(models.PersonalDetailsInput{Surname: &surname})
	if err != nil {
		t.Fatalf("PersonalDetailsUpdate() failed: %s", err)
	}

	if updateReturn.Surname != surname {
		t.Errorf("PersonalDetailsUpdate expected to return details with proper surname")
	}

	details, err := db.PersonalDetails()
	if err != nil {
		t.Fatalf("PersonalDetails() failed: %s", err)
	}

	if details.Surname != surname {
		t.Errorf("PersonalDetail.Surname: got %q, expected %q", details.Surname, surname)
	}
}

func TestPersonalDetailUpdateCountry(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd1/file.db"
		country  = "France"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	updateReturn, err := db.PersonalDetailsUpdate(models.PersonalDetailsInput{Country: &country})
	if err != nil {
		t.Fatalf("PersonalDetailsUpdate() failed: %s", err)
	}

	if updateReturn.Country != country {
		t.Errorf("PersonalDetailsUpdate expected to return details with proper country")
	}

	details, err := db.PersonalDetails()
	if err != nil {
		t.Fatalf("PersonalDetails() failed: %s", err)
	}

	if details.Country != country {
		t.Errorf("PersonalDetail.Country: got %q, expected %q", details.Country, country)
	}
}

func TestPersonalDetailUpdateBirthDate(t *testing.T) {
	var (
		fileName  = "/tmp/test_dir_i2i/pd1/file.db"
		birthdate = "20-10-1980"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	updateReturn, err := db.PersonalDetailsUpdate(models.PersonalDetailsInput{BirthDate: &birthdate})
	if err != nil {
		t.Fatalf("PersonalDetailsUpdate() failed: %s", err)
	}

	if updateReturn.BirthDate != birthdate {
		t.Errorf("PersonalDetailsUpdate expected to return details with proper birthdate")
	}

	details, err := db.PersonalDetails()
	if err != nil {
		t.Fatalf("PersonalDetails() failed: %s", err)
	}

	if details.BirthDate != birthdate {
		t.Errorf("PersonalDetail.BirthDate: got %q, expected %q", details.BirthDate, birthdate)
	}
}

func TestIdentityAdd(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd1/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newIdentity := models.IdentityInput{
		DisplayName: "private",
	}

	added, err := db.IdentityAdd(newIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	if added.DisplayName != newIdentity.DisplayName {
		t.Fatalf("IdentityAdd() expected to return %q got %q", newIdentity.DisplayName, added.DisplayName)
	}

	if len(added.ID) != identifierLength {
		t.Errorf("IdentityAdd() expected to have %d len identifier, got %d", identifierLength, len(added.ID))
	}
}

func TestIdentity(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd1/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newWorkIdentity := models.IdentityInput{
		DisplayName: "work",
	}
	newPrivateIdentity := models.IdentityInput{
		DisplayName: "private",
	}

	_, err = db.IdentityAdd(newWorkIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	_, err = db.IdentityAdd(newPrivateIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	identities, err := db.IdentityList()
	if err != nil {
		t.Fatalf("identities() failed: %s", err)
	}

	if len(identities) != 2 {
		t.Fatalf("Identity() expected to return 2 result got %d", len(identities))
	}

	if identities[0].DisplayName != newWorkIdentity.DisplayName && identities[0].DisplayName != newPrivateIdentity.DisplayName {
		t.Fatalf("Identity() expected to return %q or %q", newWorkIdentity.DisplayName, newPrivateIdentity.DisplayName)
	}

	if identities[1].DisplayName != newWorkIdentity.DisplayName && identities[1].DisplayName != newPrivateIdentity.DisplayName {
		t.Fatalf("Identity() expected to return %q or %q", newWorkIdentity.DisplayName, newPrivateIdentity.DisplayName)
	}

	if identities[0].ID == identities[1].ID {
		t.Errorf("IdentityAdd() expected IDs of identities to be unique, got %q %q", identities[0].ID, identities[1].ID)
	}
}

func TestContactAddUniqID(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd1/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newIdentity := models.IdentityInput{
		DisplayName: "private",
	}

	identity, err := db.IdentityAdd(newIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	newContact1 := models.ContactInput{
		PublicKey:   models.Key32{Key: cryptography.RandomKey32()},
		DisplayName: " Tom",
		Identity:    identity.ID,
	}

	newContact2 := models.ContactInput{
		PublicKey:   models.Key32{Key: cryptography.RandomKey32()},
		DisplayName: " Matt",
		Identity:    identity.ID,
	}

	added1, err := db.ContactAdd(newContact1)
	if err != nil {
		t.Fatalf("ContactAdd() failed: %s", err)
	}

	added2, err := db.ContactAdd(newContact2)
	if err != nil {
		t.Fatalf("ContactAdd() failed: %s", err)
	}

	if added1.ID == added2.ID {
		t.Fatalf("ContactAdd() duplicated ID: %q %q", added1.ID, added2.ID)
	}
}

func TestContactAddReturnValidation(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd1/file.db"
		name     = "Tim"
		surname  = "Cook"
		country  = "Canada"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newIdentity := models.IdentityInput{
		DisplayName: "private",
	}

	identity, err := db.IdentityAdd(newIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	newContact := models.ContactInput{
		PublicKey:    models.Key32{Key: cryptography.RandomKey32()},
		DisplayName:  "cool friend",
		Name:         &name,
		Surname:      &surname,
		Country:      &country,
		Identity:     identity.ID,
		SignatureKey: models.Key32{Key: cryptography.RandomKey32()},
	}

	addedContact, err := db.ContactAdd(newContact)
	if err != nil {
		t.Fatalf("ContactAdd() failed: %s", err)
	}

	if addedContact.Country != *newContact.Country {
		t.Errorf("ContactAdd() Country field: expected %q, got %q", *newContact.Country, addedContact.Country)
	}

	if addedContact.Name != *newContact.Name {
		t.Errorf("ContactAdd() Name field: expected %q, got %q", *newContact.Name, addedContact.Name)
	}

	if addedContact.Surname != *newContact.Surname {
		t.Errorf("ContactAdd() Surname field: expected %q, got %q", *newContact.Surname, addedContact.Surname)
	}

	if addedContact.DisplayName != newContact.DisplayName {
		t.Errorf("ContactAdd() DisplayName field: expected %q, got %q", newContact.DisplayName, addedContact.DisplayName)
	}

	if !addedContact.PublicKey.Key.Equal(newContact.PublicKey.Key) {
		t.Errorf("ContactAdd() PublicKey field: expected %q, got %q", newContact.PublicKey.Key.String(), addedContact.PublicKey.Key.String())
	}

	if !addedContact.SignatureKey.Key.Equal(newContact.SignatureKey.Key) {
		t.Errorf("ContactAdd() SignatureKey field: expected %q, got %q", newContact.SignatureKey.Key.String(), addedContact.SignatureKey.Key.String())
	}
}

func TestContactAddNoIdentity(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd1/file.db"
		identity = "non_exist"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newContact := models.ContactInput{
		PublicKey:   models.Key32{Key: cryptography.RandomKey32()},
		DisplayName: " Tom",
		Identity:    identity,
	}

	if _, err := db.ContactAdd(newContact); err == nil {
		t.Fatalf("ContactAdd() with bad identity should fail")
	}
}

func TestContactListEmpty(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd1/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newIdentity := models.IdentityInput{
		DisplayName: "private",
	}

	identity, err := db.IdentityAdd(newIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	_, err = db.ContactList(identity.ID)
	if err != nil {
		t.Fatalf("ContactList() failed: %s", err)
	}
}

func TestContactListAdded(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd1/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newIdentity := models.IdentityInput{
		DisplayName: "private",
	}

	identity, err := db.IdentityAdd(newIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	contact1 := models.ContactInput{
		DisplayName: "Tom",
		Identity:    identity.ID,
		PublicKey:   models.Key32{Key: cryptography.RandomKey32()},
	}

	contact2 := models.ContactInput{
		DisplayName: "Tom",
		Identity:    identity.ID,
		PublicKey:   models.Key32{Key: cryptography.RandomKey32()},
	}

	if _, err := db.ContactAdd(contact1); err != nil {
		t.Fatalf("ContactAdd() failed: %s", err)
	}

	if _, err := db.ContactAdd(contact2); err != nil {
		t.Fatalf("ContactAdd() failed: %s", err)
	}

	contacts, err := db.ContactList(identity.ID)
	if err != nil {
		t.Fatalf("ContactList() failed: %s", err)
	}

	if contacts == nil {
		t.Fatalf("ContactList() returned nil list")
	}

	if l := len(contacts); l != 2 {
		t.Errorf("ContactList: expected to have 2 contacts, got %d", l)
	}
}

func TestAddressAddNoIdentity(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd1/file.db"
		identity = "non_exist"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newAddress := models.AddressInput{
		DisplayName: "Home",
		Identity:    identity,
	}

	if _, err := db.AddressAdd(newAddress); err == nil {
		t.Fatalf("AddressAdd() with bad identity should fail")
	}
}

func TestAddressAddReturnValidation(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd2/file.db"
		city     = "New York"
		street   = "Long"
		country  = "Canada"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newIdentity := models.IdentityInput{
		DisplayName: "private",
	}

	identity, err := db.IdentityAdd(newIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	newAddress := models.AddressInput{
		City:        &city,
		Street:      &street,
		DisplayName: "cool address",
		Country:     &country,
		Identity:    identity.ID,
	}

	addedAddress, err := db.AddressAdd(newAddress)
	if err != nil {
		t.Fatalf("AddressAdd() failed: %s", err)
	}

	if addedAddress.Country != *newAddress.Country {
		t.Errorf("AddressAdd() Country field: expected %q, got %q", *newAddress.Country, addedAddress.Country)
	}

	if addedAddress.Street != *newAddress.Street {
		t.Errorf("AddressAdd() Street field: expected %q, got %q", *newAddress.Street, addedAddress.Street)
	}

	if addedAddress.City != *newAddress.City {
		t.Errorf("AddressAdd() City field: expected %q, got %q", *newAddress.City, addedAddress.City)
	}

	if addedAddress.DisplayName != newAddress.DisplayName {
		t.Errorf("AddressAdd() DisplayName field: expected %q, got %q", newAddress.DisplayName, addedAddress.DisplayName)
	}
}

func TestAddressListAdded(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd1/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newIdentity := models.IdentityInput{
		DisplayName: "work",
	}

	identity, err := db.IdentityAdd(newIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	address1 := models.AddressInput{
		DisplayName: "Home",
		Identity:    identity.ID,
	}

	address2 := models.AddressInput{
		DisplayName: "Work",
		Identity:    identity.ID,
	}

	if _, err := db.AddressAdd(address1); err != nil {
		t.Fatalf("AddressAdd() failed: %s", err)
	}

	if _, err := db.AddressAdd(address2); err != nil {
		t.Fatalf("AddressAdd() failed: %s", err)
	}

	addresses, err := db.AddressList(identity.ID)
	if err != nil {
		t.Fatalf("AddressList() failed: %s", err)
	}

	if addresses == nil {
		t.Fatalf("AddressList() returned nil list")
	}

	if l := len(addresses); l != 2 {
		t.Errorf("AddressList() expected to return 2 addresses, got %d", l)
	}
}

func TestPaymentCardAddNoIdentity(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd6/file.db"
		identity = "non_exist"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newPaymentCard := models.PaymentCardInput{
		DisplayName: "Home",
		Identity:    identity,
	}

	if _, err := db.PaymentCardAdd(newPaymentCard); err == nil {
		t.Fatalf("PaymentCardAdd() with bad identity should fail")
	}
}

func TestPaymentCardAddReturnValidation(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd2/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newIdentity := models.IdentityInput{
		DisplayName: "private",
	}

	identity, err := db.IdentityAdd(newIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	newPaymentCard := models.PaymentCardInput{
		Surname:      "Hanks",
		Name:         "Tom",
		Expiration:   "11-11-2012",
		SecurityCode: "123",
		Number:       "11122233344",
		Currency:     "EUR",
		DisplayName:  "private card",
		Identity:     identity.ID,
	}

	addedPaymentCard, err := db.PaymentCardAdd(newPaymentCard)
	if err != nil {
		t.Fatalf("PaymentCardAdd() failed: %s", err)
	}

	if addedPaymentCard.ID == "" {
		t.Errorf("PaymentCardAdd() ID not assigned")
	}

	if addedPaymentCard.Surname != newPaymentCard.Surname {
		t.Errorf("PaymentCardAdd() Surname field: expected %q, got %q", newPaymentCard.Surname, addedPaymentCard.Surname)
	}

	if addedPaymentCard.Name != newPaymentCard.Name {
		t.Errorf("PaymentCardAdd() Name field: expected %q, got %q", newPaymentCard.Name, addedPaymentCard.Name)
	}

	if addedPaymentCard.Expiration != newPaymentCard.Expiration {
		t.Errorf("PaymentCardAdd() Expiration field: expected %q, got %q", newPaymentCard.Expiration, addedPaymentCard.Expiration)
	}

	if addedPaymentCard.SecurityCode != newPaymentCard.SecurityCode {
		t.Errorf("PaymentCardAdd() SecurityCode field: expected %q, got %q", newPaymentCard.SecurityCode, addedPaymentCard.SecurityCode)
	}

	if addedPaymentCard.Number != newPaymentCard.Number {
		t.Errorf("PaymentCardAdd() Number field: expected %q, got %q", newPaymentCard.Number, addedPaymentCard.Number)
	}

	if addedPaymentCard.Currency != newPaymentCard.Currency {
		t.Errorf("PaymentCardAdd() Currency field: expected %q, got %q", newPaymentCard.Currency, addedPaymentCard.Currency)
	}

	if addedPaymentCard.DisplayName != newPaymentCard.DisplayName {
		t.Errorf("PaymentCardAdd() DisplayName field: expected %q, got %q", newPaymentCard.DisplayName, addedPaymentCard.DisplayName)
	}
}

func TestPaymentCardListAdded(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd1/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newIdentity := models.IdentityInput{
		DisplayName: "work",
	}

	identity, err := db.IdentityAdd(newIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	newPaymentCard1 := models.PaymentCardInput{
		Surname:      "Color",
		Name:         "Alice",
		Expiration:   "11-11-2013",
		SecurityCode: "123",
		Number:       "11122233344",
		Currency:     "USD",
		DisplayName:  "public card",
		Identity:     identity.ID,
	}

	newPaymentCard2 := models.PaymentCardInput{
		Surname:      "Hanks",
		Name:         "Tom",
		Expiration:   "11-11-2012",
		SecurityCode: "123",
		Number:       "11122233344",
		Currency:     "EUR",
		DisplayName:  "private card",
		Identity:     identity.ID,
	}

	if _, err := db.PaymentCardAdd(newPaymentCard1); err != nil {
		t.Fatalf("PaymentCardAdd() failed: %s", err)
	}

	if _, err := db.PaymentCardAdd(newPaymentCard2); err != nil {
		t.Fatalf("PaymentCardAdd() failed: %s", err)
	}

	paymentCards, err := db.PaymentCardList(identity.ID)
	if err != nil {
		t.Fatalf("PaymentCardList() failed: %s", err)
	}

	if paymentCards == nil {
		t.Fatalf("PaymentCardList() returned nil list")
	}

	if l := len(paymentCards); l != 2 {
		t.Errorf("PaymentCardList() expected to return 2 addresses, got %d", l)
	}
}

func TestIdentityDocumentAddNoIdentity(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd6/file.db"
		identity = "non_exist"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newIdentityDocument := models.IdentityDocumentInput{
		DisplayName: "personal",
		Identity:    identity,
	}

	if _, err := db.IdentityDocumentAdd(newIdentityDocument); err == nil {
		t.Fatalf("IdentityDocumentAdd() with bad identity should fail")
	}
}

func TestIdentityDocumentAddReturnValidation(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd2/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newIdentity := models.IdentityInput{
		DisplayName: "private",
	}

	identity, err := db.IdentityAdd(newIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	newIdentityDocument := models.IdentityDocumentInput{
		Surname:     "Hanks",
		Name:        "Tom",
		Expiration:  "11-11-2012",
		Number:      "11122233344",
		DisplayName: "private document",
		Country:     "Norway",
		Identity:    identity.ID,
	}

	addedDocument, err := db.IdentityDocumentAdd(newIdentityDocument)
	if err != nil {
		t.Fatalf("IdentityDocumentAdd() failed: %s", err)
	}

	if addedDocument.ID == "" {
		t.Errorf("IdentityDocumentAdd() ID not assigned")
	}

	if addedDocument.Surname != newIdentityDocument.Surname {
		t.Errorf("IdentityDocumentAdd() Surname field: expected %q, got %q", newIdentityDocument.Surname, addedDocument.Surname)
	}

	if addedDocument.Name != newIdentityDocument.Name {
		t.Errorf("IdentityDocumentAdd() Name field: expected %q, got %q", newIdentityDocument.Name, addedDocument.Name)
	}

	if addedDocument.Expiration != newIdentityDocument.Expiration {
		t.Errorf("IdentityDocumentAdd() Expiration field: expected %q, got %q", newIdentityDocument.Expiration, addedDocument.Expiration)
	}

	if addedDocument.Number != newIdentityDocument.Number {
		t.Errorf("IdentityDocumentAdd() Number field: expected %q, got %q", newIdentityDocument.Number, addedDocument.Number)
	}

	if addedDocument.Country != newIdentityDocument.Country {
		t.Errorf("IdentityDocumentAdd() Country field: expected %q, got %q", newIdentityDocument.Country, addedDocument.Country)
	}

	if addedDocument.DisplayName != newIdentityDocument.DisplayName {
		t.Errorf("IdentityDocumentAdd() DisplayName field: expected %q, got %q", newIdentityDocument.DisplayName, addedDocument.DisplayName)
	}
}

func TestIdentityDocumentListAdded(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd2/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newIdentity := models.IdentityInput{
		DisplayName: "work",
	}

	identity, err := db.IdentityAdd(newIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	newIdentityDocument1 := models.IdentityDocumentInput{
		Surname:     "Color",
		Name:        "Alice",
		Expiration:  "11-11-2013",
		Number:      "123",
		Country:     "Poland",
		DisplayName: "public card",
		Identity:    identity.ID,
	}

	newIdentityDocument2 := models.IdentityDocumentInput{
		Surname:     "Stall",
		Name:        "Tom",
		Expiration:  "11-12-2019",
		Number:      "222",
		Country:     "Spain",
		DisplayName: "private card",
		Identity:    identity.ID,
	}

	if _, err := db.IdentityDocumentAdd(newIdentityDocument1); err != nil {
		t.Fatalf("IdentityDocumentAdd() failed: %s", err)
	}

	if _, err := db.IdentityDocumentAdd(newIdentityDocument2); err != nil {
		t.Fatalf("IdentityDocumentAdd() failed: %s", err)
	}

	identityDocuments, err := db.IdentityDocumentList(identity.ID)
	if err != nil {
		t.Fatalf("IdentityDocumentList() failed: %s", err)
	}

	if identityDocuments == nil {
		t.Fatalf("IdentityDocumentList() returned nil list")
	}

	if l := len(identityDocuments); l != 2 {
		t.Errorf("IdentityDocumentList() expected to return 2 addresses, got %d", l)
	}
}

func TestIdentityDocumentDel(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd7/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newIdentity := models.IdentityInput{
		DisplayName: "delete",
	}

	identity, err := db.IdentityAdd(newIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	newIdentityDocument := models.IdentityDocumentInput{
		Surname:     "Color",
		Name:        "Alice",
		Expiration:  "11-11-2013",
		Number:      "123",
		Country:     "Poland",
		DisplayName: "public card",
		Identity:    identity.ID,
	}

	addedDocument, err := db.IdentityDocumentAdd(newIdentityDocument)
	if err != nil {
		t.Fatalf("IdentityDocumentAdd() failed: %s", err)
	}

	removedID, err := db.IdentityDocumentDel(addedDocument.ID)
	if err != nil {
		t.Fatalf("IdentityDocumentDel() failed: %s", err)
	}

	identityDocuments, err := db.IdentityDocumentList(identity.ID)
	if err != nil {
		t.Fatalf("IdentityDocumentList() failed: %s", err)
	}

	if removedID != addedDocument.ID {
		t.Errorf("IdentityDocumentDel returned %q, expected %q", removedID, addedDocument.ID)
	}

	if l := len(identityDocuments); l != 0 {
		t.Errorf("IdentityDocumentList() expected to return 0 identity documents, got %d", l)
	}
}

func TestPaymentCardDel(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd1/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newIdentity := models.IdentityInput{
		DisplayName: "work",
	}

	identity, err := db.IdentityAdd(newIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	newPaymentCard := models.PaymentCardInput{
		Surname:      "Color",
		Name:         "Alice",
		Expiration:   "11-11-2013",
		SecurityCode: "123",
		Number:       "11122233344",
		Currency:     "USD",
		DisplayName:  "public card",
		Identity:     identity.ID,
	}

	addedPaymentCard, err := db.PaymentCardAdd(newPaymentCard)
	if err != nil {
		t.Fatalf("PaymentCardAdd() failed: %s", err)
	}

	removed, err := db.PaymentCardDel(addedPaymentCard.ID)
	if err != nil {
		t.Fatalf("PaymentCardDel() failed: %s", err)
	}

	if removed != addedPaymentCard.ID {
		t.Errorf("PaymentCardDel() returned %q, expected %q", removed, addedPaymentCard)
	}

	paymentCards, err := db.PaymentCardList(identity.ID)
	if err != nil {
		t.Fatalf("PaymentCardList() failed: %s", err)
	}

	if l := len(paymentCards); l != 0 {
		t.Errorf("PaymentCardList() expected to return 0 payment cards, got %d", l)
	}
}

func TestContactDel(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd1/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newIdentity := models.IdentityInput{
		DisplayName: "private",
	}

	identity, err := db.IdentityAdd(newIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	contact := models.ContactInput{
		DisplayName: "Tom",
		Identity:    identity.ID,
		PublicKey:   models.Key32{Key: cryptography.RandomKey32()},
	}

	addedContact, err := db.ContactAdd(contact)
	if err != nil {
		t.Fatalf("ContactAdd() failed: %s", err)
	}

	removedID, err := db.ContactDel(addedContact.ID)
	if err != nil {
		t.Fatalf("ContactDel() failed: %s", err)
	}

	if removedID != addedContact.ID {
		t.Errorf("ContactDel() returned %q, expected %q", removedID, addedContact.ID)
	}

	contacts, err := db.ContactList(identity.ID)
	if err != nil {
		t.Fatalf("ContactList() failed: %s", err)
	}

	if l := len(contacts); l != 0 {
		t.Errorf("ContactList() expected to return 0 contacts, got %d", l)
	}
}

func TestAddressDel(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/d1/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newIdentity := models.IdentityInput{
		DisplayName: "work",
	}

	identity, err := db.IdentityAdd(newIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	address := models.AddressInput{
		DisplayName: "Home",
		Identity:    identity.ID,
	}

	addedAddress, err := db.AddressAdd(address)
	if err != nil {
		t.Fatalf("AddressAdd() failed: %s", err)
	}

	removedID, err := db.AddressDel(addedAddress.ID)
	if err != nil {
		t.Fatalf("AddressDel() failed: %s", err)
	}

	if removedID != addedAddress.ID {
		t.Errorf("AddressDel returned %q, expected %q", removedID, addedAddress.ID)
	}

	addresses, err := db.AddressList(identity.ID)
	if err != nil {
		t.Fatalf("AddressList() failed: %s", err)
	}

	if l := len(addresses); l != 0 {
		t.Errorf("AddressList() expected to return 0 addresses, got %d", l)
	}
}

func TestPassportAddNoIdentity(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd64/file.db"
		identity = "non_exist"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newPassport := models.PassportInput{
		DisplayName: "personal",
		Identity:    identity,
	}

	if _, err := db.PassportAdd(newPassport); err == nil {
		t.Fatalf("IdentityDocumentAdd() with bad identity should fail")
	}
}

func TestPassportAddReturnValidation(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd2/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newIdentity := models.IdentityInput{
		DisplayName: "private",
	}

	identity, err := db.IdentityAdd(newIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	newPassport := models.PassportInput{
		Surname:     "Hanks",
		Name:        "Tom",
		Expiration:  "11-11-2012",
		Number:      "11122233344",
		DisplayName: "private document",
		Country:     "Norway",
		Identity:    identity.ID,
	}

	addedPassport, err := db.PassportAdd(newPassport)
	if err != nil {
		t.Fatalf("PassportAdd() failed: %s", err)
	}

	if addedPassport.ID == "" {
		t.Errorf("PassportAdd() ID not assigned")
	}

	if addedPassport.Surname != newPassport.Surname {
		t.Errorf("PassportAdd() Surname field: expected %q, got %q", newPassport.Surname, addedPassport.Surname)
	}

	if addedPassport.Name != newPassport.Name {
		t.Errorf("PassportAdd() Name field: expected %q, got %q", newPassport.Name, addedPassport.Name)
	}

	if addedPassport.Expiration != newPassport.Expiration {
		t.Errorf("PassportAdd() Expiration field: expected %q, got %q", newPassport.Expiration, addedPassport.Expiration)
	}

	if addedPassport.Number != newPassport.Number {
		t.Errorf("PassportAdd() Number field: expected %q, got %q", newPassport.Number, addedPassport.Number)
	}

	if newPassport.Country != newPassport.Country {
		t.Errorf("PassportAdd() Country field: expected %q, got %q", newPassport.Country, addedPassport.Country)
	}

	if newPassport.DisplayName != newPassport.DisplayName {
		t.Errorf("PassportAdd() DisplayName field: expected %q, got %q", newPassport.DisplayName, newPassport.DisplayName)
	}
}

func TestPassportListAdded(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pdf2/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newIdentity := models.IdentityInput{
		DisplayName: "work",
	}

	identity, err := db.IdentityAdd(newIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	newPassport1 := models.PassportInput{
		Surname:     "Color",
		Name:        "Alice",
		Expiration:  "11-11-2013",
		Number:      "123",
		Country:     "Poland",
		DisplayName: "public card",
		Identity:    identity.ID,
	}

	newPassport2 := models.PassportInput{
		Surname:     "Stall",
		Name:        "Tom",
		Expiration:  "11-12-2019",
		Number:      "222",
		Country:     "Spain",
		DisplayName: "private card",
		Identity:    identity.ID,
	}

	if _, err := db.PassportAdd(newPassport1); err != nil {
		t.Fatalf("PassportAdd() failed: %s", err)
	}

	if _, err := db.PassportAdd(newPassport2); err != nil {
		t.Fatalf("PassportAdd() failed: %s", err)
	}

	passports, err := db.PassportList(identity.ID)
	if err != nil {
		t.Fatalf("PassportList() failed: %s", err)
	}

	if passports == nil {
		t.Fatalf("PassportList() returned nil list")
	}

	if l := len(passports); l != 2 {
		t.Errorf("PassportList() expected to return 2 addresses, got %d", l)
	}
}

func TestPassportDel(t *testing.T) {
	var (
		fileName = "/tmp/test_dir_i2i/pd7x/file.db"
	)

	wallet, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain() failed: %s", err)
	}

	db, err := LoadDatabase(fileName, wallet)
	if err != nil {
		t.Fatalf("LoadDatabase(%q) failed: %s", fileName, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %s", err)
		}

		if err := os.RemoveAll("/tmp/test_dir_i2i"); err != nil {
			t.Errorf("failed to clean after test: %s", err)
		}
	}()

	newIdentity := models.IdentityInput{
		DisplayName: "delete",
	}

	identity, err := db.IdentityAdd(newIdentity)
	if err != nil {
		t.Fatalf("IdentityAdd() failed: %s", err)
	}

	newPassport := models.PassportInput{
		Surname:     "Color",
		Name:        "Alice",
		Expiration:  "11-11-2013",
		Number:      "123",
		Country:     "Poland",
		DisplayName: "public card",
		Identity:    identity.ID,
	}

	addedPassport, err := db.PassportAdd(newPassport)
	if err != nil {
		t.Fatalf("PassportAdd() failed: %s", err)
	}

	removedID, err := db.PassportDel(addedPassport.ID)
	if err != nil {
		t.Fatalf("PassportDel() failed: %s", err)
	}

	passports, err := db.PassportList(identity.ID)
	if err != nil {
		t.Fatalf("PassportList() failed: %s", err)
	}

	if removedID != addedPassport.ID {
		t.Errorf("PassportDel returned %q, expected %q", removedID, addedPassport.ID)
	}

	if l := len(passports); l != 0 {
		t.Errorf("PassportList() expected to return 0 identity documents, got %d", l)
	}
}

func bucketExist(db *bolt.DB, bucketName string, t *testing.T) {
	err := db.Update(func(tx *bolt.Tx) error {
		if tx.Bucket([]byte(bucketName)) == nil {
			return ErrBucketNotFound(bucketName)
		}
		return nil
	})

	if err != nil {
		t.Errorf("bucketExist: %s", err)
	}
}
