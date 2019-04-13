package generator

import "testing"

func TestNames(t *testing.T) {
	for i := range names {
		if names[i] == "" {
			t.Errorf("name should not be empty")
		}
	}
}
