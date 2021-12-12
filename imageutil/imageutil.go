package imageutil

import (
	"image"
	"image/color"
	"math"
)

func CalculateDistance(i1 image.Image, i2 image.Image) float64 {
	r := 0.0
	xmax := i1.Bounds().Dx()
	if xmax > i2.Bounds().Dx() {
		xmax = i2.Bounds().Dx()
	}
	ymax := i1.Bounds().Dy()
	if ymax > i2.Bounds().Dy() {
		ymax = i2.Bounds().Dy()
	}
	for x := 0; x < xmax; x++ {
		for y := 0; y < xmax; y++ {
			c1r, c1b, c1g, _ := i1.At(x, y).RGBA()
			c2r, c2b, c2g, _ := i2.At(x, y).RGBA()
			r += (math.Abs(float64(c1r-c2r))/255.0 +
				math.Abs(float64(c1b-c2b))/255.0 +
				math.Abs(float64(c1g-c2g))/255.0) / 3.0
		}
	}
	return r
}

type CopyOnRead struct {
	Copy *image.RGBA
	image.Image
}

func (cow *CopyOnRead) At(x, y int) color.Color {
	if cow.Copy == nil {
		cow.Copy = image.NewRGBA(cow.Image.Bounds())
	}
	c := cow.Image.At(x, y)
	cow.Copy.Set(x, y, c)
	return c
}
