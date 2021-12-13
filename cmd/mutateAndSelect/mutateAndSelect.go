package main

import (
	"encoding/csv"
	"fmt"
	"github.com/agnivade/levenshtein"
	"image"
	"image-formula-find/dna1"
	"image-formula-find/imageutil"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"math"
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
	const generations = 100
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
		const mutations = 8
		var children = make([]*imageutil.Individual, 0, len(lastGeneration)*mutations+len(lastGeneration)*len(lastGeneration)+childrenCount+1)

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
					Lineage:         p.DNA,
				})
			}
		}

		//for i := 0; i < len(lastGeneration)*len(lastGeneration); i++ {
		//	p1 := lastGeneration[i%len(lastGeneration)]
		//	p2 := lastGeneration[i/len(lastGeneration)]
		if len(lastGeneration) > 4 {
			p1 := lastGeneration[int(rand.Int31n(int32(len(lastGeneration))))]
			p2 := lastGeneration[int(rand.Int31n(int32(len(lastGeneration))))]
			dna := dna1.Breed(p1.DNA, p2.DNA)
			if _, ok := seen[dna]; ok {
				continue
			}
			seen[dna] = struct{}{}
			if dna1.Valid(dna) {
				if dna != p1.DNA && dna != p2.DNA && p1.Lineage != p2.Lineage {
					children = append(children, &imageutil.Individual{
						DNA: dna,
						Parent: []*imageutil.Individual{
							p1, p2,
						},
						Lineage:         dna,
						FirstGeneration: generation,
					})
				}
			}
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
				Lineage:         dna,
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
		for len(lastGeneration) < childrenCount && len(children) > 0 {
			child := children[0]
			children = children[1:]

			minDistance := math.MaxInt
			for _, lg := range lastGeneration {
				m := levenshtein.ComputeDistance(lg.DNA, child.DNA)
				if m < minDistance {
					minDistance = m
					if minDistance < 10 {
						break
					}
				}
			}

			if minDistance < 10 {
				continue
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
	fin, err := os.Open("in4.png")
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
