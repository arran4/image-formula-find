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
	"sort"
	"sync"
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

type Sorter struct {
	children []string
	scores   []float64
}

func (s *Sorter) Len() int {
	return len(s.children)
}

func (s *Sorter) Less(i, j int) bool {
	return s.scores[i] < s.scores[j]
}

func (s *Sorter) Swap(i, j int) {
	s.scores[i], s.scores[j] = s.scores[j], s.scores[i]
	s.children[i], s.children[j] = s.children[j], s.children[i]
}

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	rand.Seed(time.Now().UnixNano())
	const logGenerations = 10
	const generations = 1000
	const childrenCount = 10

	srcimg := LoadImage()

	plotSize := srcimg.Bounds()

	destimg := image.NewRGBA(image.Rect(0, 0, plotSize.Dx()*childrenCount, plotSize.Dy()*logGenerations))

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
		row = append(row,
			fmt.Sprintf("C%d Dna", i+1),
			fmt.Sprintf("C%d Formula Red", i+1),
			fmt.Sprintf("C%d Formula Blue", i+1),
			fmt.Sprintf("C%d Formula Green", i+1),
			fmt.Sprintf("C%d Distance", i+1))
	}
	headerSize := len(row)
	csvw.Write(row)
	for generation := 0; generation < generations; generation++ {
		log.Printf("Generation %d", generation+1)
		scores := make([]float64, len(children), len(children))
		row = make([]string, headerSize, headerSize)
		wg := sync.WaitGroup{}
		for fi := range children {
			wg.Add(1)
			go func(i int, child string) {
				defer wg.Done()
				rd, bd, gd := dna1.SplitString3(child)
				rf := dna1.ParseFunction(rd)
				bf := dna1.ParseFunction(bd)
				gf := dna1.ParseFunction(gd)
				d := &drawer1.Drawer{
					RedFormula:   rf,
					BlueFormula:  bf,
					GreenFormula: gf,
				}
				childImage := image.NewRGBA(plotSize.Bounds())
				if (generation % (generations / logGenerations)) > 0 {
					draw.Draw(childImage, plotSize, d, image.Pt(0, 0), draw.Src)
				} else {
					cow := &CopyOnRead{Copy: childImage, Image: d}
					draw.Draw(destimg, plotSize.Add(image.Pt(plotSize.Dx()*i, plotSize.Dy()*(generation/(generations/logGenerations)))), cow, image.Pt(0, 0), draw.Src)
				}
				distance := imageutil.CalculateDistance(srcimg, childImage)
				copy(row[i*5:(i+1)*5], []string{
					child, rf.String(), bf.String(), gf.String(), fmt.Sprintf("%0.2f", distance),
				})
				scores[i] = distance
			}(fi, children[fi])
		}
		wg.Wait()
		csvw.Write(row)

		sort.Sort(&Sorter{
			children: children,
			scores:   scores,
		})

		for i := 0; i < 3; i++ {
			for ii := 0; ii < 3; ii++ {
				if i == ii {
					children[i*3+ii] = dna1.Mutate(children[i])
				} else {
					children[i*3+1] = dna1.Breed(children[i], children[ii])
				}
			}
		}
		children[9] = dna1.RndStr(50)
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
