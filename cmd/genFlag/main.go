package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func main() {
	// Mauritius flag: Red, Blue, Yellow, Green
	// 4 horizontal stripes.
	// Flag dims: 60x40
	// Border: 5px
	flagW, flagH := 60, 40
	border := 5
	width := flagW + border*2
	height := flagH + border*2

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill white background (border)
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Pt(0, 0), draw.Src)

	red := color.RGBA{234, 40, 57, 255}
	blue := color.RGBA{26, 32, 109, 255}
	yellow := color.RGBA{255, 213, 0, 255}
	green := color.RGBA{0, 165, 81, 255}

	stripeHeight := flagH / 4

	for y := 0; y < flagH; y++ {
		var c color.RGBA
		if y < stripeHeight {
			c = red
		} else if y < stripeHeight*2 {
			c = blue
		} else if y < stripeHeight*3 {
			c = yellow
		} else {
			c = green
		}
		for x := 0; x < flagW; x++ {
			img.Set(x+border, y+border, c)
		}
	}

	f, _ := os.Create("flag.png")
	png.Encode(f, img)
	f.Close()
}
