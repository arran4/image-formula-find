package drawer1

import (
	"image"
	image_formula_find "image-formula-find"
	"image/color"
	"sync"
)

type Drawer struct {
	RedFormula   *image_formula_find.Function
	BlueFormula  *image_formula_find.Function
	GreenFormula *image_formula_find.Function
}

func (d *Drawer) Convert(c color.Color) color.Color {
	return c
}

func (d *Drawer) ColorModel() color.Model {
	return d
}

func (d *Drawer) Bounds() image.Rectangle {
	return image.Rectangle{image.Point{-1e9, -1e9}, image.Point{1e9, 1e9}}
}

func (d *Drawer) At(x, y int) color.Color {
	wg := sync.WaitGroup{}
	wg.Add(3)
	var rr float64
	var gr float64
	var br float64
	go func() {
		defer wg.Done()
		rr, _, _ = d.RedFormula.Evaluate(float64(x), float64(y), 0)
	}()
	go func() {
		defer wg.Done()
		br, _, _ = d.BlueFormula.Evaluate(float64(x), float64(y), 0)
	}()
	go func() {
		defer wg.Done()
		gr, _, _ = d.GreenFormula.Evaluate(float64(x), float64(y), 0)
	}()
	wg.Wait()
	return color.RGBA{
		R: uint8(rr),
		G: uint8(gr),
		B: uint8(br),
		A: 255,
	}
}
