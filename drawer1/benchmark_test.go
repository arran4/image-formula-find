package drawer1

import (
	"image-formula-find"
	"testing"
)

func BenchmarkDrawerAt_Simple(b *testing.B) {
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

func BenchmarkDrawerAt_Heavy(b *testing.B) {
	// Heavy formula: Construct a deep tree
	vX := &image_formula_find.Var{Var: "X"}
	vY := &image_formula_find.Var{Var: "Y"}

	var expr image_formula_find.Expression = &image_formula_find.Plus{LHS: vX, RHS: vY}
	// Nest operations to create computational load
	for i := 0; i < 20; i++ {
		expr = &image_formula_find.Multiply{
			LHS: expr,
			RHS: image_formula_find.NewSingleFunction(
				"Sin",
				&image_formula_find.Plus{
					LHS: vY,
					RHS: &image_formula_find.Const{Value: float64(i)},
				},
			),
		}
	}

	formula := &image_formula_find.Function{
		Equals: &image_formula_find.Equals{
			LHS: &image_formula_find.Const{Value: 0},
			RHS: expr,
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
