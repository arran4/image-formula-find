package dna1

import (
	"fmt"
	"image-formula-find"
	"testing"
)

func TestParseConstValues(t *testing.T) {
	// chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	// Op mapping: i % 17
	// 0: ParseConst (exp 0)  ('A')
	// 1-11: Ops
	// 12: ParseConst (exp 1) ('M', '/')
	// 13: ParseConst (exp 2)
	// 14: ParseConst (exp -1) ('w')
	// 15: ParseConst (exp -2) ('g')
	// 16: MakeConst ('Q')

	testCases := []struct {
		input string
		desc  string
		expected float64
		rest string
	}{
		// 'M' is index 12. 12%17 = 12 -> ParseConst (exp 1).
		// Arg empty. r=0. val=0.
		{"M", "ParseConst 'M' (exp 1) empty arg", 0.0, ""},

		// "ABA": 'A' (0->exp 0). 'B'(1), 'A'(0). r=64. val=64.
		{"ABA", "ParseConst 'A' (exp 0, x1) val 64", 64.0, ""},

		// "/BA": '/' (63->12->exp 1). 'B'(1), 'A'(0). r=64. val=640.
		{"/BA", "ParseConst '/' (exp 1, x10) val 640", 640.0, ""},

		// "gBA": 'g' (32->15->exp -2). r=64. val=0.64.
		{"gBA", "ParseConst 'g' (exp -2, x0.01) val 0.64", 0.64, ""},

		// "wBA": 'w' (48->14->exp -1). r=64. val=6.4.
		{"wBA", "ParseConst 'w' (exp -1, x0.1) val 6.4", 6.4, ""},

		// "Q": 'Q' (16->16->MakeConst). val=16.
		{"Q", "MakeConst 'Q' (val 16)", 16.0, ""},

		// "QBA": 'Q' -> MakeConst(16). Rest "BA".
		{"QBA", "MakeConst 'Q' (val 16) with rest", 16.0, "BA"},
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

		if rest != tc.rest {
			t.Errorf("%s: expected rest %q, got %q", tc.desc, tc.rest, rest)
		}
	}
}
