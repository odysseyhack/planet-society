package generator

import "testing"

func TestGnames(t *testing.T) {
	if gnames() == nil {
		t.Errorf("gnames returned nil")
	}
}

func TestGsurnames(t *testing.T) {
	if gsurnames() == nil {
		t.Errorf("gsurnames returned nil")
	}
}
func TestGcountries(t *testing.T) {
	if gcountries() == nil {
		t.Errorf("gcountries returned nil")
	}
}
func TestGstreets(t *testing.T) {
	if gstreets() == nil {
		t.Errorf("gstreets returned nil")
	}
}
func TestGidentityNames(t *testing.T) {
	if gidentityNames() == nil {
		t.Errorf("gidentityNames returned nil")
	}
}
func TestGaddressName(t *testing.T) {
	if gaddressName() == nil {
		t.Errorf("gaddressName returned nil")
	}
}
func TestCurrencies(t *testing.T) {
	if currencies() == nil {
		t.Errorf("currencies returned nil")
	}
}
func TestGrequestReasons(t *testing.T) {
	if grequestReasons() == nil {
		t.Errorf("grequestReasons returned nil")
	}
}
