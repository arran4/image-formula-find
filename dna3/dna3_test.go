package dna3

import (
	"strings"
	"testing"
)

func TestResolve(t *testing.T) {
	// Helper to make string of length 30 filled with 'A' (0)
	padding := strings.Repeat("A", 29) // 29 As + 1 char at start makes 30

	tests := []struct {
		name     string
		dna      string
		index    int
		expected string // String representation of expression
	}{
		{
			name:     "Base Constant A",
			dna:      "A",
			index:    0,
			expected: "0",
		},
		{
			name:     "Base Constant B",
			dna:      "B", // B (index 1) -> 25
			index:    0,
			expected: "25",
		},
		{
			name:     "Additive Layer",
			dna:      "B" + padding + "B", // Index 0: B(25), Index 30: B(25)
			index:    0,
			expected: "25 + 25",
		},
		{
			name:     "Op Layer Sin",
			dna:      "B" + padding + "2", // Index 0: B(25), Index 30: 2(Sin - 54) - Wait, '2' is index 54? No.
			// chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
			// '2' is at index 54. MapOp(54) is Sin.
			index:    0,
			expected: "Sin(25)",
		},
		{
			name:     "Op Layer Cos",
			dna:      "B" + padding + "3", // Index 0: B(25), Index 30: 3(Cos - 55)
			// '3' is at index 55. MapOp(55) is Cos.
			index:    0,
			expected: "Cos(25)",
		},
		{
			name:     "Multi Layer",
			dna:      "B" + padding + "B" + padding + "2", // 0:B, 30:B, 60:2(Sin)
			index:    0,
			expected: "Sin(25 + 25)",
		},
		{
			name:     "Empty Layer Ignored",
			dna:      "B" + padding + "A" + padding + "B", // 0:B, 30:A, 60:B
			index:    0,
			expected: "25 + 25", // A is skipped
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := Resolve(tt.index, tt.dna)
			got := expr.String()
			if got != tt.expected {
				t.Errorf("Resolve() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseChannel(t *testing.T) {
	// Test basic structure
	// Need enough DNA for at least one param?
	// ParseChannel loops 6 times.
	// If DNA is short, it uses 0.

	dna := "BBBBB" // 5 params, all 1 -> 25.
	// Term 0 (even): X * 25^25 + 25*25 + 25
	// Other terms: 0
	// Total: 0 + (X * 25^25 + 25*25 + 25) + ...

	expr := ParseChannel(dna)
	s := expr.String()

	// Just check if it contains expected parts
	if !strings.Contains(s, "X * 25 ^ 25") {
		t.Errorf("Expected X * 25 ^ 25, got %s", s)
	}
	if !strings.Contains(s, "25 * 25") {
		t.Errorf("Expected 25 * 25, got %s", s)
	}
}
