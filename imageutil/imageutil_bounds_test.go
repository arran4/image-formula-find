package imageutil

import (
	"fmt"
	"image"
	"image/color"
	"testing"
)

type StrictBoundsImage struct {
	Rect     image.Rectangle
	Accessed map[image.Point]bool
}

func (i *StrictBoundsImage) ColorModel() color.Model { return color.RGBAModel }
func (i *StrictBoundsImage) Bounds() image.Rectangle { return i.Rect }
func (i *StrictBoundsImage) At(x, y int) color.Color {
	pt := image.Point{x, y}
	if !pt.In(i.Rect) {
		panic(fmt.Sprintf("out of bounds access at %v", pt))
	}
	if i.Accessed == nil {
		i.Accessed = make(map[image.Point]bool)
	}
	i.Accessed[pt] = true
	return color.RGBA{0, 0, 0, 255}
}

func TestPanicOnOOB(t *testing.T) {
	// Case 1: Width > Height (xmax > ymax)
	// xmax = 10, ymax = 5
	// Loop goes to y < xmax (10). Should access (0, 5) and panic.
	w, h := 10, 5
	img1 := &StrictBoundsImage{Rect: image.Rect(0, 0, w, h)}
	img2 := &StrictBoundsImage{Rect: image.Rect(0, 0, w, h)}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panicked: %v", r)
		}
	}()

	CalculateDistance(img1, img2)
}

func TestCorrectCoverage(t *testing.T) {
	// Case 2: Width < Height (xmax < ymax)
	// xmax = 5, ymax = 10
	// Loop goes to y < xmax (5). Should stop early and miss pixels.
	w, h := 5, 10
	img1 := &StrictBoundsImage{Rect: image.Rect(0, 0, w, h), Accessed: make(map[image.Point]bool)}
	img2 := &StrictBoundsImage{Rect: image.Rect(0, 0, w, h), Accessed: make(map[image.Point]bool)}

	CalculateDistance(img1, img2)

	// Check if all pixels were accessed
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			pt := image.Point{x, y}
			if !img1.Accessed[pt] {
				t.Errorf("Pixel %v was not accessed", pt)
			}
		}
	}
}
