package drawer1

import (
	"image"
	image_formula_find "image-formula-find"
	"image/color"
	"runtime"
	"sync"
)

type Drawer struct {
	RedFormula    *image_formula_find.Function
	BlueFormula   *image_formula_find.Function
	GreenFormula  *image_formula_find.Function
	Width, Height int

	once  sync.Once
	cache *image.RGBA
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
	// If within bounds, use cached parallel pre-computation
	if d.Width > 0 && d.Height > 0 && x >= 0 && x < d.Width && y >= 0 && y < d.Height {
		d.once.Do(func() {
			d.cache = image.NewRGBA(image.Rect(0, 0, d.Width, d.Height))
			numWorkers := runtime.NumCPU()
			rowsPerWorker := (d.Height + numWorkers - 1) / numWorkers
			var wg sync.WaitGroup
			for i := 0; i < numWorkers; i++ {
				startY := i * rowsPerWorker
				endY := startY + rowsPerWorker
				if startY >= d.Height {
					break
				}
				if endY > d.Height {
					endY = d.Height
				}
				wg.Add(1)
				go func(y1, y2 int) {
					defer wg.Done()
					for cy := y1; cy < y2; cy++ {
						sy := (float64(cy)/float64(d.Height))*20.0 - 10.0
						for cx := 0; cx < d.Width; cx++ {
							sx := (float64(cx)/float64(d.Width))*20.0 - 10.0
							rr, _, _ := d.RedFormula.Evaluate(sx, sy, 0)
							br, _, _ := d.BlueFormula.Evaluate(sx, sy, 0)
							gr, _, _ := d.GreenFormula.Evaluate(sx, sy, 0)
							d.cache.SetRGBA(cx, cy, color.RGBA{
								R: uint8(rr),
								G: uint8(gr),
								B: uint8(br),
								A: 255,
							})
						}
					}
				}(startY, endY)
			}
			wg.Wait()
		})
		return d.cache.At(x, y)
	}

	// Fallback for out-of-bounds queries (or invalid dimensions)
	var rr float64
	var gr float64
	var br float64

	sx := float64(x)
	sy := float64(y)
	if d.Width > 0 && d.Height > 0 {
		sx = (float64(x)/float64(d.Width))*20.0 - 10.0
		sy = (float64(y)/float64(d.Height))*20.0 - 10.0
	}

	rr, _, _ = d.RedFormula.Evaluate(sx, sy, 0)
	br, _, _ = d.BlueFormula.Evaluate(sx, sy, 0)
	gr, _, _ = d.GreenFormula.Evaluate(sx, sy, 0)
	return color.RGBA{
		R: uint8(rr),
		G: uint8(gr),
		B: uint8(br),
		A: 255,
	}
}
