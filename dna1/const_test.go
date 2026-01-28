package dna1

import (
	"fmt"
	"image-formula-find"
	"testing"
)

func TestParseConstValues(t *testing.T) {
	// chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	// chars[0] = 'A'
	// chars[16] = 'Q'
	// chars[32] = 'g'
	// chars[48] = 'w'
	// chars[63] = '/'

	testCases := []struct {
		input string
		desc  string
		expected float64
	}{
		// 'M' is index 12. Not an operator. Should be MakeConst.
		{"M", "MakeConst (index 12)", 12.0},

		// "ABA": 'A' exp 0. 'B'(1), 'A'(0). r=64. val=64.
		{"ABA", "ParseConst 'A' (exp 0, x1) val 64", 64.0},

		// "/BA": '/' exp -2. 'B'(1), 'A'(0). r=64. val=0.64.
		{"/BA", "ParseConst '/' (exp -2, x0.01) val 64", 0.64},

		// "gBA": 'g' exp 2. r=64. val=6400.
		{"gBA", "ParseConst 'g' (exp 2, x100) val 64", 6400.0},

		// "wBA": 'w' exp -1. r=64. val=6.4.
		{"wBA", "ParseConst 'w' (exp -1, x0.1) val 64", 6.4},

		// "QBA": 'Q' exp 1. r=64. val=640.
		{"QBA", "ParseConst 'Q' (exp 1, x10) val 64", 640.0},
	}

	for _, tc := range testCases {
		rest, expr := ParseExpression(tc.input)
		fmt.Printf("%s -> %v (rest: %q)\n", tc.desc, expr, rest)

		c, ok := expr.(*image_formula_find.Const)
		if !ok {
			t.Errorf("%s: expected Const, got %T (%v)", tc.desc, expr, expr)
			continue
		}

		// allow small float error
		if diff := c.Value - tc.expected; diff > 0.0001 || diff < -0.0001 {
			t.Errorf("%s: expected %g, got %g", tc.desc, tc.expected, c.Value)
		}

		if tc.input == "M" {
			if rest != "" {
				t.Errorf("%s: expected empty rest, got %q", tc.desc, rest)
			}
		} else {
			if rest != "" {
				t.Errorf("%s: expected empty rest, got %q", tc.desc, rest)
			}
		}
	}
}
