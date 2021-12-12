package main

import (
	"fmt"
	"image"
	"image-formula-find/dna1"
	"image-formula-find/imageutil"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	childrenCount = 10
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	rand.Seed(time.Now().UnixNano())
	worker := NewWorker()
	go worker.Work()
	ebiten.SetWindowSize(640*2, 480*3)
	ebiten.SetWindowTitle("Watch Generator")
	if err := ebiten.RunGame(worker); err != nil {
		log.Fatal(err)
	}
}

type WorkerDetails struct {
	sync.RWMutex
	LastGeneration []*imageutil.Individual
	srcimg         image.Image
	plotSize       image.Rectangle
	generation     int
	esrcimg        *ebiten.Image
	Winners        []*imageutil.Individual
}

func (worker *WorkerDetails) Update() error {
	return nil
}

func (worker *WorkerDetails) Draw(screen *ebiten.Image) {
	worker.RWMutex.RLock()
	defer worker.RWMutex.RUnlock()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Generation %d", worker.generation))
	op := &ebiten.DrawImageOptions{}
	const offsetY = 20
	op.GeoM.Translate(0, offsetY)
	if worker.esrcimg == nil {
		worker.esrcimg = ebiten.NewImageFromImage(worker.srcimg)
	}
	screen.DrawImage(worker.esrcimg, op)
	for i, ind := range worker.LastGeneration {
		op := &ebiten.DrawImageOptions{}
		tx := 0 //(worker.srcimg.Bounds().Dx()+10)
		ty := (worker.srcimg.Bounds().Dy() + 10) * (i + 1)
		op.GeoM.Translate(float64(tx), float64(ty))
		i := ind.Image()
		if i != nil {
			screen.DrawImage(ebiten.NewImageFromImage(i), op)
		}
		tx += (worker.srcimg.Bounds().Dx() + 10)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score %f", ind.Score), tx, ty)
		ty += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("DNA %s", ind.DNA), tx, ty)
		ty += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Red   = %s", ind.Rf.String()), tx, ty)
		ty += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Green = %s", ind.Gf.String()), tx, ty)
		ty += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Blue  = %s", ind.Bf.String()), tx, ty)
	}
}

func (worker *WorkerDetails) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (worker *WorkerDetails) PlotSize() image.Rectangle {
	return worker.plotSize
}

func (worker *WorkerDetails) SourceImage() image.Image {
	return worker.srcimg
}

func NewWorker() *WorkerDetails {
	worker := &WorkerDetails{
		srcimg: LoadImage(),
	}
	worker.plotSize = worker.srcimg.Bounds()
	return worker
}

func (worker *WorkerDetails) Work() {
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
	for generation := 0; ; generation++ {
		log.Printf("Generation %d", generation+1)
		worker.RLock()
		lastGeneration := worker.LastGeneration
		worker.RUnlock()

		const mutations = 5
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
				})
			}
		}

		for i := 0; i < len(lastGeneration)*len(lastGeneration); i++ {
			p1 := lastGeneration[i%10]
			p2 := lastGeneration[i/10]
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
				child.Calculate(worker)
			}(fi, children[fi])
		}
		wg.Wait()

		sort.Sort(sort.Reverse(&imageutil.Sorter{
			Children: children,
		}))

		worker.Lock()
		lastGeneration = children[:10]
		children = children[10:]
		for i := range lastGeneration {
			if len(children) == 0 {
				break
			}
			child := lastGeneration[i]
			if generation-child.FirstGeneration > 10 {
				lastGeneration[i] = children[0]
				worker.Winners = append(worker.Winners, child)
				children = children[0:]
			}
		}
		sort.Sort(sort.Reverse(&imageutil.Sorter{
			Children: lastGeneration,
		}))
		sort.Sort(sort.Reverse(&imageutil.Sorter{
			Children: worker.Winners,
		}))
		if len(worker.Winners) > 100 {
			worker.Winners = worker.Winners[:100]
		}
		worker.LastGeneration = lastGeneration
		worker.generation = generation
		worker.Unlock()
	}
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
