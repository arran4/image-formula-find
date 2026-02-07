package main

import (
	"image"
	image_formula_find "image-formula-find"
	"image-formula-find/drawer1"
	"image/png"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	i := image.NewRGBA(image.Rect(0, 0, 100, 100))
	d := &drawer1.Drawer{
		RedFormula:   image_formula_find.ParseFunction("x = y + 1"),
		BlueFormula:  image_formula_find.ParseFunction("x = 2 * y + 3"),
		GreenFormula: image_formula_find.ParseFunction("x = 4 * y + 5"),
	}
	d.Render(i)
	f, err := os.Create("out.png")
	if err != nil {
		log.Panicf("Error: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()
	if err := png.Encode(f, i); err != nil {
		log.Panicf("Error: %v", err)
	}
	log.Printf("Done")
}
