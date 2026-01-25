package main

import (
	"flag"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"image-formula-find/dna1"
)

func main() {
	var inputPath string
	var outputPath string
	var generations int
	var steps int

	flag.StringVar(&inputPath, "input", "in5.png", "Path to input image")
	flag.StringVar(&outputPath, "output", "evolution.gif", "Path to output GIF")
	flag.IntVar(&generations, "generations", 1000, "Number of generations")
	flag.IntVar(&steps, "steps", 10, "Number of steps (frames in GIF)")
	flag.Parse()

	log.SetFlags(log.Flags() | log.Lshortfile)

	srcimg := LoadImage(inputPath)
	plotSize := srcimg.Bounds()
	width := plotSize.Dx()
	height := plotSize.Dy()

	// Using BasicRequired implementation from dna1
	worker := &dna1.BasicRequired{
		R: plotSize,
		I: srcimg,
	}

	const childrenCount = 10
	var lastGeneration []*dna1.Individual
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

	var frames []*image.Paletted
	var delays []int

	// Calculate interval to capture frames
    // We want to capture exactly 'steps' frames if possible.
    stepInterval := generations / steps
    if stepInterval < 1 {
        stepInterval = 1
    }

	log.Printf("Starting evolution. Generations: %d, Steps: %d, Interval: %d", generations, steps, stepInterval)

	for generation := 0; generation < generations; generation++ {
		lastGeneration = dna1.GenerationProcess(worker, lastGeneration, generation, newDNA)

		// Capture frame if it's time
		if (generation+1)%stepInterval == 0 || generation == generations-1 {
			// Ensure we don't capture more frames than steps unless it's the very last one and we want to be sure to include it.
            // But user asked for 10 steps.

			log.Printf("Capturing frame at generation %d", generation+1)

            // lastGeneration is sorted by score (best first)
            best := lastGeneration[0]
            evolvedImg := best.Image()

            // Create composite image: [Evolved | Target]
            compositeRect := image.Rect(0, 0, width*2, height)
            compositeImg := image.NewRGBA(compositeRect)

            // Draw Evolved
            draw.Draw(compositeImg, image.Rect(0, 0, width, height), evolvedImg, image.Pt(0, 0), draw.Src)
            // Draw Target
            draw.Draw(compositeImg, image.Rect(width, 0, width*2, height), srcimg, image.Pt(0, 0), draw.Src)

            // Convert to Paletted for GIF
            // Use Plan9 palette which is standard. For better quality we could generate palette from image,
            // but for simplicity and speed Plan9 is often used in basic Go GIF examples.
            // Or better: image/gif handles quantization if we use EncodeAll with options?
            // Actually gif.EncodeAll expects []*image.Paletted. So we must quantize.

            palettedImg := image.NewPaletted(compositeRect, palette.Plan9)
            draw.FloydSteinberg.Draw(palettedImg, compositeRect, compositeImg, image.Pt(0, 0))

            frames = append(frames, palettedImg)
            delays = append(delays, 10) // 100ms delay
		}
	}

	fout, err := os.Create(outputPath)
	if err != nil {
		log.Panicf("Error: %v", err)
	}
	defer fout.Close()

	g := &gif.GIF{
		Image: frames,
		Delay: delays,
	}
	if err := gif.EncodeAll(fout, g); err != nil {
		log.Panicf("Error: %v", err)
	}
	log.Printf("Saved GIF to %s with %d frames", outputPath, len(frames))
}

func LoadImage(path string) image.Image {
	fin, err := os.Open(path)
	if err != nil {
		log.Panicf("Error opening image %s: %v", path, err)
	}
	defer fin.Close()
	i, _, err := image.Decode(fin)
	if err != nil {
		log.Panicf("Error decoding image %s: %v", path, err)
	}
	return i
}
