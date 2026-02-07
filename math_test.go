package image_formula_find

import (
	"testing"
)

func TestParseFunction_InvalidInput(t *testing.T) {
	cases := []string{
		"x =",
		"(",
		")",
		"x = y +",
	}

	for _, c := range cases {
		_, err := ParseFunction(c)
		if err == nil {
			t.Errorf("Expected error for invalid formula '%s', got nil", c)
		}
	}
}

func TestParseFunction_ValidInput(t *testing.T) {
	f, err := ParseFunction("x = y + 1")
	if err != nil {
		t.Errorf("Expected no error for valid formula, got %v", err)
	}
	if f == nil {
		t.Error("Expected function, got nil")
	}
}
