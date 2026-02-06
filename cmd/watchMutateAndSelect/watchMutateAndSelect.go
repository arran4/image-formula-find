package main

import (
	"fmt"
	"image"
	"image-formula-find/dna1"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	LastGeneration []*dna1.Individual
	srcimg         image.Image
	plotSize       image.Rectangle
	generation     int
	esrcimg        *ebiten.Image
	Winners        []*dna1.Individual
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

		lastGeneration = dna1.GenerationProcess(worker, lastGeneration, generation, newDNA)

		worker.Lock()
		sort.Sort((&dna1.Sorter{
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
	fin, err := os.Open("in5.png")
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
