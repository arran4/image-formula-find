package drawer1

import (
	"image-formula-find"
	"testing"
)

func BenchmarkDrawerAt(b *testing.B) {
	// Simple formula: X + Y
	formula := &image_formula_find.Function{
		Equals: &image_formula_find.Equals{
			LHS: &image_formula_find.Const{Value: 0},
			RHS: &image_formula_find.Plus{
				LHS: &image_formula_find.Var{Var: "X"},
				RHS: &image_formula_find.Var{Var: "Y"},
			},
		},
	}

	d := &Drawer{
		RedFormula:   formula,
		GreenFormula: formula,
		BlueFormula:  formula,
		Width:        1000,
		Height:       1000,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.At(i%1000, (i/1000)%1000)
	}
}
