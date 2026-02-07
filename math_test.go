package image_formula_find

import (
	"testing"
)

func TestEvaluateCorrectness(t *testing.T) {
	exprStr := "y = x * 2"
	f, err := ParseFunction(exprStr)
	if err != nil {
		t.Fatalf("Failed to parse function: %v", err)
	}
	if f == nil {
		t.Fatal("Failed to parse function")
	}

	// For y=2, x=1, t=0 -> RHS - LHS = (1*2) - 2 = 0
	w, _, err := f.Evaluate(1, 2, 0)
	if err != nil {
		t.Fatal(err)
	}
	if w != 0 {
		t.Errorf("Expected weight 0, got %f", w)
	}

	// x=2, y=2 -> (2*2) - 2 = 2
	w, _, _ = f.Evaluate(2, 2, 0)
	if w != 2 {
		t.Errorf("Expected weight 2, got %f", w)
	}
}
 
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
