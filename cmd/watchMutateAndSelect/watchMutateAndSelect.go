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
	ebiten.SetWindowSize(640, 480)
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
	for _, ind := range worker.LastGeneration {
		op.GeoM.Translate(float64(worker.srcimg.Bounds().Dx()+10), 0)
		i := ind.Image()
		if i != nil {
			screen.DrawImage(ebiten.NewImageFromImage(i), op)
		}
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
	for generation := 0; ; generation++ {
		log.Printf("Generation %d", generation+1)
		worker.RLock()
		lastGeneration := worker.LastGeneration
		worker.RUnlock()

		var children = make([]*imageutil.Individual, 0, len(lastGeneration)*len(lastGeneration)+childrenCount+1)

		children = append(children, lastGeneration...)

		for i := 0; i < len(lastGeneration)*len(lastGeneration); i++ {
			p1 := lastGeneration[i%10]
			p2 := lastGeneration[i/10]
			children = append(children, &imageutil.Individual{
				DNA: dna1.Breed(p1.DNA, p2.DNA),
				Parent: []*imageutil.Individual{
					p1, p2,
				},
			})
		}

		for {
			children = append(children, &imageutil.Individual{
				DNA: dna1.RndStr(50),
			})
			if len(children) > childrenCount {
				break
			}
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
