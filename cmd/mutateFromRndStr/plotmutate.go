package main

import (
	"encoding/csv"
	"fmt"
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
	plotSize := image.Rect(0, 0, 100, 100)
	const generations = 10
	const childrenCount = 10
	img := image.NewRGBA(image.Rect(0, 0, plotSize.Dx()*childrenCount, plotSize.Dy()*generations))

	fcsv, err := os.Create("out.csv")
	if err != nil {
		log.Panicf("Error: %v", err)
	}
	defer func() {
		if err := fcsv.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()
	csvw := csv.NewWriter(fcsv)
	defer csvw.Flush()
	var row []string
	var children []string
	for i := 0; i < childrenCount; i++ {
		children = append(children, dna1.RndStr(50))
		row = append(row, fmt.Sprintf("C%d Dna", i+1), fmt.Sprintf("C%d Formula Red", i+1), fmt.Sprintf("C%d Formula Blue", i+1), fmt.Sprintf("C%d Formula Green", i+1))
	}
	if err := csvw.Write(row); err != nil {
		log.Panicf("Error writing csv: %v", err)
	}
	for generation := 0; generation < generations; generation++ {
		log.Printf("Generation %d", generation+1)
		row = []string{}
		for i, child := range children {
			row = append(row, child)
			rd, bd, gd := dna1.SplitString3(child)
			rf := dna1.ParseFunction(rd)
			bf := dna1.ParseFunction(bd)
			gf := dna1.ParseFunction(gd)
			row = append(row, rf.String(), bf.String(), gf.String())
			d := &drawer1.Drawer{
				RedFormula:   rf,
				BlueFormula:  bf,
				GreenFormula: gf,
			}
			d.Render(img.SubImage(plotSize.Add(image.Pt(plotSize.Dx()*i, plotSize.Dy()*generation))).(draw.Image))
		}
		if err := csvw.Write(row); err != nil {
			log.Panicf("Error writing csv: %v", err)
		}

		for i := 0; i < childrenCount; i++ {
			children[i] = dna1.Mutate(children[i])
		}
	}
	f, err := os.Create("out.png")
	if err != nil {
		log.Panicf("Error: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()
	if err := png.Encode(f, img); err != nil {
		log.Panicf("Error: %v", err)
	}
	log.Printf("Done")
}
