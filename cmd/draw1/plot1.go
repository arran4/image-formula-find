package main

import (
	"image"
	image_formula_find "image-formula-find"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
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

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	i := image.NewRGBA(image.Rect(0, 0, 100, 100))
	draw.Draw(i, i.Rect, &Drawer{
		RedFormula:   image_formula_find.ParseFunction("x = y + 1"),
		BlueFormula:  image_formula_find.ParseFunction("x = 2 * y + 3"),
		GreenFormula: image_formula_find.ParseFunction("x = 4 * y + 5"),
	}, image.Pt(0, 0), draw.Src)
	f, err := os.Create("out.png")
	if err != nil {
		log.Panicf("Error: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, i); err != nil {
		log.Panicf("Error: %v", err)
	}
	log.Printf("Done")
}
