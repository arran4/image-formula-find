package worker

import (
	"image"
	"image-formula-find/dna6"
	"log"
	"sort"
	"sync"
)

type Worker struct {
	sync.RWMutex
	LastGeneration []*dna6.Individual
	SrcImg         image.Image
	PlotSizeRect   image.Rectangle
	Generation     int
	Winners        []*dna6.Individual
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
			dna := dna6.RndStr(50)
			if !dna6.Valid(dna) {
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

		// Worker implements dna6.Required
		lastGeneration = dna6.GenerationProcess(worker, lastGeneration, generation, newDNA)

		worker.Lock()
		sort.Sort((&dna6.Sorter{
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
