package protocol

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/odysseyhack/planet-society/protocol/cryptography"
	"github.com/odysseyhack/planet-society/protocol/database"
	"github.com/phob0s-pl/generator"
)

func TestFillTransaction(t *testing.T) {

	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create temporary dir: %s", err)
	}
	defer os.RemoveAll(dir)

	keychain, err := cryptography.OneShotKeychain()
	if err != nil {
		t.Fatalf("OneShotKeychain failed: %s", err)
	}

	db, err := database.LoadDatabase(filepath.Join(dir, "db"), keychain)
	if err != nil {
		t.Fatalf("OneShotKeychain failed: %s", err)
	}
	defer db.Close()

	generator := generator.NewGenerator()
	if err := generator.Generate(db); err != nil {
		t.Fatalf("Generate failed: %s", err)
	}

	transaction, err := fillTransaction(db)
	if err != nil {
		t.Errorf("fillTransaction failed %s", err)
	}

	if transaction == nil {
		t.Errorf("fillTransaction returned nil transaction")
	}
}
