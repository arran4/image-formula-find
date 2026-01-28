package drawer1

import (
	"image-formula-find"
	"testing"
)

func TestDrawerScaling(t *testing.T) {
	// Formula: X + 10
	// We use LHS=0, RHS=(X+10).
	// Var("X") evaluates to scaled X.

	// Construct X + 10 manually to avoid parsing dependency if possible,
	// or just use parser if available. Parser is in image_formula_find package
	// but mostly internal or generated. ParseFunction is exported.

	// Let's use manual construction to be safe and simple.

	xPlus10 := &image_formula_find.Function{
		Equals: &image_formula_find.Equals{
			LHS: &image_formula_find.Const{Value: 0},
			RHS: &image_formula_find.Plus{
				LHS: &image_formula_find.Var{Var: "X"},
				RHS: &image_formula_find.Const{Value: 10},
			},
		},
	}

	// Same for Y
	yPlus10 := &image_formula_find.Function{
		Equals: &image_formula_find.Equals{
			LHS: &image_formula_find.Const{Value: 0},
			RHS: &image_formula_find.Plus{
				LHS: &image_formula_find.Var{Var: "Y"},
				RHS: &image_formula_find.Const{Value: 10},
			},
		},
	}

	zero := &image_formula_find.Function{
		Equals: &image_formula_find.Equals{
			LHS: &image_formula_find.Const{Value: 0},
			RHS: &image_formula_find.Const{Value: 0},
		},
	}

	d := &Drawer{
		RedFormula:   xPlus10,
		GreenFormula: yPlus10,
		BlueFormula:  zero,
		Width:        100,
		Height:       100,
	}

	// Test X scaling (Red channel)
	// x=0 -> scaled -10. +10 -> 0.
	c := d.At(0, 0)
	r, _, _, _ := c.RGBA() // Returns alpha-premultiplied values in [0, 65535]

	if r != 0 {
		t.Errorf("At(0,0) Red (X): expected 0, got %d", r)
	}

	// x=50 -> scaled 0. +10 -> 10.
	// 10 * 0x101 = 2570.
	c = d.At(50, 0)
	r, _, _, _ = c.RGBA()
	expected10 := uint32(10) * 0x101
	if r != expected10 {
		t.Errorf("At(50,0) Red (X): expected %d (val 10), got %d", expected10, r)
	}

	// x=100 -> scaled 10. +10 -> 20.
	c = d.At(100, 0)
	r, _, _, _ = c.RGBA()
	expected20 := uint32(20) * 0x101
	if r != expected20 {
		t.Errorf("At(100,0) Red (X): expected %d (val 20), got %d", expected20, r)
	}

	// Test Y scaling (Green channel)
	// y=0 -> scaled -10. +10 -> 0.
	c = d.At(0, 0)
	_, g, _, _ := c.RGBA()
	if g != 0 {
		t.Errorf("At(0,0) Green (Y): expected 0, got %d", g)
	}

	// y=100 -> scaled 10. +10 -> 20.
	c = d.At(0, 100)
	_, g, _, _ = c.RGBA()
	if g != expected20 {
		t.Errorf("At(0,100) Green (Y): expected %d (val 20), got %d", expected20, g)
	}
}
