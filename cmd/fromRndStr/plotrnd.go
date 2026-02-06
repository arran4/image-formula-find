package main

import (
	"image"
	"image-formula-find/dna1"
	"image-formula-find/drawer1"
	"image/draw"
	"image/png"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	i := image.NewRGBA(image.Rect(0, 0, 100, 100))
	rf, bf, gf := dna1.SplitString3(dna1.RndStr(18))
	d := &drawer1.Drawer{
		RedFormula:   dna1.ParseFunction(rf),
		BlueFormula:  dna1.ParseFunction(bf),
		GreenFormula: dna1.ParseFunction(gf),
	}
	draw.Draw(i, i.Rect, d, image.Pt(0, 0), draw.Src)
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
