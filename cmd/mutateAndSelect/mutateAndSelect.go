package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"image-formula-find"
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
	children []*Individual
}

func (s *Sorter) Len() int {
	return len(s.children)
}

func (s *Sorter) Less(i, j int) bool {
	return s.children[i].Score < s.children[j].Score
}

func (s *Sorter) Swap(i, j int) {
	s.children[i], s.children[j] = s.children[j], s.children[i]
}

type Individual struct {
	DNA    string
	Parent []*Individual
	Score  float64
	rf     *image_formula_find.Function
	bf     *image_formula_find.Function
	gf     *image_formula_find.Function
	i      draw.Image
	d      *drawer1.Drawer
}

func (i *Individual) Calculate(plotSize image.Rectangle, srcimg image.Image) {
	rd, bd, gd := dna1.SplitString3(i.DNA)
	i.rf = dna1.ParseFunction(rd)
	i.bf = dna1.ParseFunction(bd)
	i.gf = dna1.ParseFunction(gd)
	i.d = &drawer1.Drawer{
		RedFormula:   i.rf,
		BlueFormula:  i.bf,
		GreenFormula: i.gf,
	}
	i.i = image.NewRGBA(plotSize.Bounds())
	draw.Draw(i.i, plotSize, i.d, image.Pt(0, 0), draw.Src)
	i.Score = imageutil.CalculateDistance(srcimg, i.i)
}

func (i *Individual) CsvRow() []string {
	return []string{
		i.DNA, i.rf.String(), i.bf.String(), i.gf.String(), fmt.Sprintf("%0.2f", i.Score),
	}
}

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	rand.Seed(time.Now().UnixNano())
	const logGenerations = 10
	const generations = 10
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
	var lastGeneration []*Individual
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
	for generation := 0; generation < generations; generation++ {
		log.Printf("Generation %d", generation+1)
		var children = make([]*Individual, 0, len(lastGeneration)*len(lastGeneration)+childrenCount+1)

		children = append(children, lastGeneration...)

		for i := 0; i < len(lastGeneration)*len(lastGeneration); i++ {
			p1 := lastGeneration[i%10]
			p2 := lastGeneration[i/10]
			children = append(children, &Individual{
				DNA: dna1.Breed(p1.DNA, p2.DNA),
				Parent: []*Individual{
					p1, p2,
				},
			})
		}

		for {
			children = append(children, &Individual{
				DNA: dna1.RndStr(50),
			})
			if len(children) > childrenCount {
				break
			}
		}

		row = make([]string, 0, headerSize)
		wg := sync.WaitGroup{}
		for fi := range children {
			wg.Add(1)
			go func(i int, child *Individual) {
				defer wg.Done()
				child.Calculate(plotSize, srcimg)
			}(fi, children[fi])
		}
		wg.Wait()

		sort.Sort(sort.Reverse(&Sorter{
			children: children,
		}))

		if (generation % (generations / logGenerations)) == 0 {
			for i, child := range children[:10] {
				draw.Draw(destimg, plotSize.Add(image.Pt(plotSize.Dx()*i, plotSize.Dy()*(generation/(generations/logGenerations)))), child.i, image.Pt(0, 0), draw.Src)
				row = append(row, child.CsvRow()...)
			}
		}
		csvw.Write(row)

		lastGeneration = children[:10]
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
