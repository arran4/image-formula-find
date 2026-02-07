//go:build noebiten

package main

import (
	"fmt"
	"image"
	"image/png"
	"image-formula-find/worker"
	"log"
	"os"
	"time"
)

type Imager interface {
	Image() image.Image
}

func run(w *worker.Worker) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			w.RLock()
			gen := w.Generation
			var best Imager
			if len(w.Winners) > 0 {
				best = w.Winners[0]
			} else if len(w.LastGeneration) > 0 {
				best = w.LastGeneration[0]
			}
			w.RUnlock()

			if best == nil {
				continue
			}

			img := best.Image()
			if img == nil {
				continue
			}

			filename := fmt.Sprintf("gen_%d.png", gen)
			f, err := os.Create(filename)
			if err != nil {
				log.Printf("Failed to create file: %v", err)
				continue
			}
			if err := png.Encode(f, img); err != nil {
				log.Printf("Failed to encode png: %v", err)
			}
			f.Close()
			log.Printf("Saved %s", filename)
		}
	}
}
