package drawer1

import (
	"image"
	image_formula_find "image-formula-find"
	"image/color"
)

type Drawer struct {
	RedFormula    *image_formula_find.Function
	BlueFormula   *image_formula_find.Function
	GreenFormula  *image_formula_find.Function
	Width, Height int
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
	sx := float64(x)
	sy := float64(y)
	if d.Width > 0 && d.Height > 0 {
		sx = (float64(x)/float64(d.Width))*20.0 - 10.0
		sy = (float64(y)/float64(d.Height))*20.0 - 10.0
	}

	rr, _, _ := d.RedFormula.Evaluate(sx, sy, 0)
	br, _, _ := d.BlueFormula.Evaluate(sx, sy, 0)
	gr, _, _ := d.GreenFormula.Evaluate(sx, sy, 0)

	return color.RGBA{
		R: uint8(rr),
		G: uint8(gr),
		B: uint8(br),
		A: 255,
	}
}
