package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"image-formula-find/dna1"
	"image-formula-find/drawer1"
	"image-formula-find/imageutil"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"
)

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

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	rand.Seed(time.Now().UnixNano())
	const generations = 10
	const childrenCount = 10

	srcimg := LoadImage()

	plotSize := srcimg.Bounds()

	destimg := image.NewRGBA(image.Rect(0, 0, plotSize.Dx()*childrenCount, plotSize.Dy()*generations))

	fcsv, err := os.Create("out.csv")
	if err != nil {
		log.Panicf("Error: %v", err)
	}
	defer fcsv.Close()
	csvw := csv.NewWriter(fcsv)
	defer csvw.Flush()
	var row []string
	var children []string
	for i := 0; i < childrenCount; i++ {
		children = append(children, dna1.RndStr(50))
		row = append(row, fmt.Sprintf("C%d Dna", i+1), fmt.Sprintf("C%d Formula Red", i+1), fmt.Sprintf("C%d Formula Blue", i+1), fmt.Sprintf("C%d Formula Green", i+1), fmt.Sprintf("C%d Distance", i+1))
	}
	csvw.Write(row)
	for generation := 0; generation < generations; generation++ {
		log.Printf("Generation %d", generation+1)
		scores := []float64{}
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
			cow := &CopyOnRead{Copy: image.NewRGBA(plotSize.Bounds()), Image: d}
			draw.Draw(destimg, plotSize.Add(image.Pt(plotSize.Dx()*i, plotSize.Dy()*generation)), cow, image.Pt(0, 0), draw.Src)
			distance := imageutil.CalculateDistance(srcimg, cow.Copy)
			row = append(row, fmt.Sprintf("%0.2f", distance))
			scores = append(scores, distance)
		}
		csvw.Write(row)

		for i := 0; i < childrenCount; i++ {
			children[i] = dna1.Mutate(children[i])
		}
	}
	fout, err := os.Create("out.png")
	if err != nil {
		log.Panicf("Error: %v", err)
	}
	defer fout.Close()
	if err := png.Encode(fout, destimg); err != nil {
		log.Panicf("Error: %v", err)
	}
	log.Printf("Done")
}

func LoadImage() image.Image {
	fin, err := os.Open("in.png")
	if err != nil {
		log.Panicf("Error: %v", err)
	}
	defer fin.Close()
	i, _, err := image.Decode(fin)
	if err != nil {
		log.Panicf("error: %v", err)
	}
	return i
}
