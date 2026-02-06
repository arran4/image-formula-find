package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"image-formula-find/dna1"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	const logGenerations = 10
	const generations = 1000
	const childrenCount = 10

	srcimg := LoadImage()

	plotSize := srcimg.Bounds()

	destimg := image.NewRGBA(image.Rect(0, 0, plotSize.Dx()*(childrenCount+1), plotSize.Dy()*logGenerations))

	draw.Draw(destimg, srcimg.Bounds().Add(image.Pt(plotSize.Dx()*(childrenCount), plotSize.Dy()*(logGenerations-1))), srcimg, image.Pt(0, 0), draw.Src)

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
	var lastGeneration []*dna1.Individual
	for i := 0; i < childrenCount; i++ {
		row = append(row,
			fmt.Sprintf("C%d Dna", i+1),
			fmt.Sprintf("C%d Formula Red", i+1),
			fmt.Sprintf("C%d Formula Blue", i+1),
			fmt.Sprintf("C%d Formula Green", i+1),
			fmt.Sprintf("C%d Distance", i+1))
	}
	headerSize := len(row)
	if err := csvw.Write(row); err != nil {
		log.Panicf("Error writing csv: %v", err)
	}
	newDNA := make(chan string, 100)
	go func() {
		for {
			dna := dna1.RndStr(50)
			if !dna1.Valid(dna) {
				continue
			}
			newDNA <- dna
		}
	}()
	worker := &dna1.BasicRequired{
		R: plotSize,
		I: srcimg,
	}
	for generation := 0; generation < generations; generation++ {
		log.Printf("Generation %d", generation+1)

		lastGeneration = dna1.GenerationProcess(worker, lastGeneration, generation, newDNA)

		if (generation % (generations / logGenerations)) == 0 {
			row = make([]string, 0, headerSize)
			for i, child := range lastGeneration {
				draw.Draw(destimg, plotSize.Add(image.Pt(plotSize.Dx()*i, plotSize.Dy()*(generation/(generations/logGenerations)))), child.Image(), image.Pt(0, 0), draw.Src)
				row = append(row, child.CsvRow()...)
			}
			if err := csvw.Write(row); err != nil {
				log.Panicf("Error writing csv: %v", err)
			}
		}

	}
	fout, err := os.Create("out.png")
	if err != nil {
		log.Panicf("Error: %v", err)
	}
	defer func() {
		if err := fout.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()
	if err := png.Encode(fout, destimg); err != nil {
		log.Panicf("Error: %v", err)
	}
	log.Printf("Done")
}

func LoadImage() image.Image {
	fin, err := os.Open("in5.png")
	if err != nil {
		log.Panicf("Error: %v", err)
	}
	defer func() {
		if err := fin.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()
	i, _, err := image.Decode(fin)
	if err != nil {
		log.Panicf("error: %v", err)
	}
	return i
}
