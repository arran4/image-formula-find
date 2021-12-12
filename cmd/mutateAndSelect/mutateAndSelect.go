package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"image-formula-find/dna1"
	"image-formula-find/imageutil"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	rand.Seed(time.Now().UnixNano())
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
	defer fcsv.Close()
	csvw := csv.NewWriter(fcsv)
	defer csvw.Flush()
	var row []string
	var lastGeneration []*imageutil.Individual
	for i := 0; i < childrenCount; i++ {
		row = append(row,
			fmt.Sprintf("C%d Dna", i+1),
			fmt.Sprintf("C%d Formula Red", i+1),
			fmt.Sprintf("C%d Formula Blue", i+1),
			fmt.Sprintf("C%d Formula Green", i+1),
			fmt.Sprintf("C%d Distance", i+1))
	}
	headerSize := len(row)
	csvw.Write(row)
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
	for generation := 0; generation < generations; generation++ {
		log.Printf("Generation %d", generation+1)
		const mutations = 5
		var children = make([]*imageutil.Individual, 0, len(lastGeneration)*mutations+len(lastGeneration)*len(lastGeneration)+childrenCount+1)

		children = append(children, lastGeneration...)

		seen := map[string]struct{}{}

		for _, p := range lastGeneration {
			dna := p.DNA
			if _, ok := seen[dna]; ok {
				continue
			}
			seen[dna] = struct{}{}
			children = append(children, p)
		}

		for _, p := range lastGeneration {
			for i := 0; i <= mutations; i++ {
				dna := p.DNA
				for m := 0; m < i; m++ {
					dna = dna1.Mutate(dna)
				}
				if _, ok := seen[dna]; ok {
					continue
				}
				seen[dna] = struct{}{}
				if !dna1.Valid(dna) {
					continue
				}
				children = append(children, &imageutil.Individual{
					DNA:             dna,
					Parent:          []*imageutil.Individual{p},
					FirstGeneration: generation,
				})
			}
		}

		for i := 0; i < len(lastGeneration)*len(lastGeneration); i++ {
			p1 := lastGeneration[i%len(lastGeneration)]
			p2 := lastGeneration[i/len(lastGeneration)]
			dna := dna1.Breed(p1.DNA, p2.DNA)
			if _, ok := seen[dna]; ok {
				continue
			}
			seen[dna] = struct{}{}
			if !dna1.Valid(dna) {
				continue
			}
			children = append(children, &imageutil.Individual{
				DNA: dna,
				Parent: []*imageutil.Individual{
					p1, p2,
				},
				FirstGeneration: generation,
			})
		}

		for len(children) < childrenCount {
			dna := <-newDNA
			if _, ok := seen[dna]; ok {
				continue
			}
			seen[dna] = struct{}{}
			if !dna1.Valid(dna) {
				continue
			}
			children = append(children, &imageutil.Individual{
				DNA:             dna,
				FirstGeneration: generation,
			})
		}

		wg := sync.WaitGroup{}
		for fi := range children {
			wg.Add(1)
			go func(i int, child *imageutil.Individual) {
				defer wg.Done()
				child.Calculate(&imageutil.BasicRequired{
					R: plotSize,
					I: srcimg,
				})
			}(fi, children[fi])
		}
		wg.Wait()

		sort.Sort(sort.Reverse(&imageutil.Sorter{
			Children: children,
		}))

		lastGeneration = make([]*imageutil.Individual, 0, childrenCount)
		lineages := map[string]int{}
		for len(lastGeneration) < childrenCount && len(children) > 0 {
			child := children[0]
			children = children[1:]
			ph := ""
			for _, p := range child.Parent {
				ph += p.DNA
			}
			if len(ph) > 0 {
				if v, ok := lineages[ph]; ok && v > 3 {
					continue
				}
				lineages[ph]++
			}
			lastGeneration = append(lastGeneration, child)
		}
		sort.Sort(sort.Reverse(&imageutil.Sorter{
			Children: lastGeneration,
		}))

		if (generation % (generations / logGenerations)) == 0 {
			row = make([]string, 0, headerSize)
			for i, child := range lastGeneration {
				draw.Draw(destimg, plotSize.Add(image.Pt(plotSize.Dx()*i, plotSize.Dy()*(generation/(generations/logGenerations)))), child.Image(), image.Pt(0, 0), draw.Src)
				row = append(row, child.CsvRow()...)
			}
			csvw.Write(row)
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
