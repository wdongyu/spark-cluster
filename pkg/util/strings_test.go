package util

import (
	"testing"
)

func TestRandRunes(t *testing.T) {
	str := RandRunes(8)
	if len(str) != 8 {
		t.Errorf("Expected length of rune %v, actual %v", 8, len(str))
	}

	for _, ch := range str {
		if !((ch >= rune('a') && ch <= rune('z')) || (ch >= rune('0') && ch <= rune('9'))) {
			t.Errorf("Found unexpected character %v", ch)
			break
		}
	}
}

func TestRandStringRunes(t *testing.T) {
	str := RandStringRunes(8)
	if len(str) != 8 {
		t.Errorf("Expected length of rune %v, actual %v", 8, len(str))
	}

	for _, ch := range str {
		if !(ch >= rune('a') && ch <= rune('z')) {
			t.Errorf("Found unexpected character %v", ch)
			break
		}
	}
}
