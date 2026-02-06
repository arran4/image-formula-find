package worker

import (
	"image"
	"image-formula-find/dna1"
	"log"
	"sort"
	"sync"
)

type Worker struct {
	sync.RWMutex
	LastGeneration []*dna1.Individual
	SrcImg         image.Image
	PlotSizeRect   image.Rectangle
	Generation     int
	Winners        []*dna1.Individual
}

func NewWorker(img image.Image) *Worker {
	return &Worker{
		SrcImg:       img,
		PlotSizeRect: img.Bounds(),
	}
}

func (w *Worker) PlotSize() image.Rectangle {
	return w.PlotSizeRect
}

func (w *Worker) SourceImage() image.Image {
	return w.SrcImg
}

func (worker *Worker) Work() {
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

		// Worker implements dna1.Required
		lastGeneration = dna1.GenerationProcess(worker, lastGeneration, generation, newDNA)

		worker.Lock()
		sort.Sort((&dna1.Sorter{
			Children: worker.Winners,
		}))
		if len(worker.Winners) > 100 {
			worker.Winners = worker.Winners[:100]
		}
		worker.LastGeneration = lastGeneration
		worker.Generation = generation
		worker.Unlock()
	}
}
