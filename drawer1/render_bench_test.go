package drawer1

import (
	"image"
	"image-formula-find"
	"image/draw"
	"testing"
)

func BenchmarkDraw_Heavy(b *testing.B) {
	// Heavy formula setup
	vX := &image_formula_find.Var{Var: "X"}
	vY := &image_formula_find.Var{Var: "Y"}
	var expr image_formula_find.Expression = &image_formula_find.Plus{LHS: vX, RHS: vY}
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

	width := 200
	height := 200
	d := &Drawer{
		RedFormula:   formula,
		GreenFormula: formula,
		BlueFormula:  formula,
		Width:        width,
		Height:       height,
	}
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		draw.Draw(dst, dst.Bounds(), d, image.Pt(0, 0), draw.Src)
	}
}

func BenchmarkRender_Heavy(b *testing.B) {
	// Same setup
	vX := &image_formula_find.Var{Var: "X"}
	vY := &image_formula_find.Var{Var: "Y"}
	var expr image_formula_find.Expression = &image_formula_find.Plus{LHS: vX, RHS: vY}
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

	width := 200
	height := 200
	d := &Drawer{
		RedFormula:   formula,
		GreenFormula: formula,
		BlueFormula:  formula,
		Width:        width,
		Height:       height,
	}
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Render(dst)
	}
}
