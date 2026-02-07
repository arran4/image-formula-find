package main

import (
	"image"
	"image-formula-find/worker"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	i := LoadImage()
	w := worker.NewWorker(i)
	go w.Work()
	run(w)
}

func LoadImage() image.Image {
	fin, err := os.Open("in5.png")
	if err != nil {
		log.Panicf("Error: %v", err)
	}
	defer func() {
		if err := fin.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()
	i, _, err := image.Decode(fin)
	if err != nil {
		log.Panicf("error: %v", err)
	}
	return i
}
