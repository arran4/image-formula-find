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
	rf, err := image_formula_find.ParseFunction("x = y + 1")
	if err != nil {
		log.Fatalf("Invalid red formula: %v", err)
	}
	bf, err := image_formula_find.ParseFunction("x = 2 * y + 3")
	if err != nil {
		log.Fatalf("Invalid blue formula: %v", err)
	}
	gf, err := image_formula_find.ParseFunction("x = 4 * y + 5")
	if err != nil {
		log.Fatalf("Invalid green formula: %v", err)
	}
	d := &drawer1.Drawer{
		RedFormula:   rf,
		BlueFormula:  bf,
		GreenFormula: gf,
	}
	d.Render(i)
	log.Printf("Red: %s", d.RedFormula.String())
	log.Printf("Blue: %s", d.BlueFormula.String())
	log.Printf("Green: %s", d.GreenFormula.String())
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
