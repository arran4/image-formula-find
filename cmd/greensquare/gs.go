package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
)

func main() {
	i := image.NewRGBA(image.Rect(0, 0, 100, 100))
	draw.Draw(i, i.Rect, image.NewUniform(color.RGBA{G: 255, A: 255}), image.Pt(0, 0), draw.Src)

	fout, err := os.Create("gs1.png")
	if err != nil {
		log.Panicf("Error: %v", err)
	}
	defer fout.Close()
	if err := png.Encode(fout, i); err != nil {
		log.Panicf("Error: %v", err)
	}
	log.Printf("Done")

}
