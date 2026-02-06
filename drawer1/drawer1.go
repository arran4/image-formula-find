package drawer1

import (
	"image"
	image_formula_find "image-formula-find"
	"image/color"
	"image/draw"
	"runtime"
	"sync"
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

// Render draws the formula to the destination image in parallel.
// It assumes the destination bounds map 1:1 to the Drawer's coordinate space (0,0 to Width,Height).
func (d *Drawer) Render(dst draw.Image) {
	bounds := dst.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	minX := bounds.Min.X
	minY := bounds.Min.Y

	numWorkers := runtime.NumCPU()
	if numWorkers < 1 {
		numWorkers = 1
	}

	// We split the work into horizontal bands
	rowsPerWorker := (height + numWorkers - 1) / numWorkers

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		startY := i * rowsPerWorker
		endY := (i + 1) * rowsPerWorker
		if endY > height {
			endY = height
		}
		// If startY >= height (can happen if more workers than rows), skip
		if startY >= height {
			wg.Done()
			continue
		}

		go func(y0, y1 int) {
			defer wg.Done()
			for y := y0; y < y1; y++ {
				// Calculate sy (scaled Y)
				sy := float64(y)
				if d.Height > 0 {
					sy = (float64(y) / float64(d.Height)) * 20.0 - 10.0
				}

				for x := 0; x < width; x++ {
					// Calculate sx (scaled X)
					sx := float64(x)
					if d.Width > 0 {
						sx = (float64(x) / float64(d.Width)) * 20.0 - 10.0
					}

					// Evaluate formulas
					// Note: Evaluate is assumed thread-safe (pure function)
					rr, _, _ := d.RedFormula.Evaluate(sx, sy, 0)
					br, _, _ := d.BlueFormula.Evaluate(sx, sy, 0)
					gr, _, _ := d.GreenFormula.Evaluate(sx, sy, 0)

					dst.Set(minX+x, minY+y, color.RGBA{
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
}
