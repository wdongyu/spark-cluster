package util

import "testing"

func TestIsExist(t *testing.T) {
	p := "/fake"
	if IsExist(p) {
		t.Errorf("Expected non exist, actual exist")
	}
}
