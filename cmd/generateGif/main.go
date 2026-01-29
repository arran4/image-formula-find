package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/arran4/golang-wordwrap"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

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
	imgWidth := plotSize.Dx()
	imgHeight := plotSize.Dy()

	// Layout configuration
	const (
		padding       = 10
		labelHeight   = 20
		formulaHeight = 150
		borderWidth   = 2
	)

	// Vertical Layout:
	// Top: Target
	// Middle: Evolution
	// Bottom: Formula

	canvasWidth := padding*2 + imgWidth
	if canvasWidth < 400 { canvasWidth = 400 } // Ensure width for text

	// Height = Padding + Label + Target + Padding + Label + Evolution + Padding + Formula + Padding
	canvasHeight := padding + labelHeight + imgHeight + padding + labelHeight + imgHeight + padding + formulaHeight + padding

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
	stepInterval := generations / steps
	if stepInterval < 1 {
		stepInterval = 1
	}

	log.Printf("Starting evolution. Generations: %d, Steps: %d, Interval: %d", generations, steps, stepInterval)

	for generation := 0; generation < generations; generation++ {
		lastGeneration = dna1.GenerationProcess(worker, lastGeneration, generation, newDNA)

		// Capture frame if it's time
		if (generation+1)%stepInterval == 0 || generation == generations-1 {
			log.Printf("Capturing frame at generation %d", generation+1)

			best := lastGeneration[0]
			evolvedImg := best.Image()

			// Create composite image
			compositeRect := image.Rect(0, 0, canvasWidth, canvasHeight)
			compositeImg := image.NewRGBA(compositeRect)

			// Fill background with white
			draw.Draw(compositeImg, compositeRect, &image.Uniform{color.White}, image.Pt(0, 0), draw.Src)

			// Helper to draw border
			drawBorder := func(r image.Rectangle, c color.Color) {
				// Top
				draw.Draw(compositeImg, image.Rect(r.Min.X-borderWidth, r.Min.Y-borderWidth, r.Max.X+borderWidth, r.Min.Y), &image.Uniform{c}, image.Pt(0, 0), draw.Src)
				// Bottom
				draw.Draw(compositeImg, image.Rect(r.Min.X-borderWidth, r.Max.Y, r.Max.X+borderWidth, r.Max.Y+borderWidth), &image.Uniform{c}, image.Pt(0, 0), draw.Src)
				// Left
				draw.Draw(compositeImg, image.Rect(r.Min.X-borderWidth, r.Min.Y, r.Min.X, r.Max.Y), &image.Uniform{c}, image.Pt(0, 0), draw.Src)
				// Right
				draw.Draw(compositeImg, image.Rect(r.Max.X, r.Min.Y, r.Max.X+borderWidth, r.Max.Y), &image.Uniform{c}, image.Pt(0, 0), draw.Src)
			}

			// Center X for images
			imgX := (canvasWidth - imgWidth) / 2

			// Draw Target (Top)
			addLabel(compositeImg, padding, padding+labelHeight-5, "Target")
			targetRect := image.Rect(imgX, padding+labelHeight, imgX+imgWidth, padding+labelHeight+imgHeight)
			drawBorder(targetRect, color.Black)
			draw.Draw(compositeImg, targetRect, srcimg, image.Pt(0, 0), draw.Src)

			// Draw Evolution (Middle)
			evoY := padding + labelHeight + imgHeight + padding
			addLabel(compositeImg, padding, evoY+labelHeight-5, "Evolution")
			evolvedRect := image.Rect(imgX, evoY+labelHeight, imgX+imgWidth, evoY+labelHeight+imgHeight)
			drawBorder(evolvedRect, color.Black)
			draw.Draw(compositeImg, evolvedRect, evolvedImg, image.Pt(0, 0), draw.Src)

			// Draw Formula
			formulaY := evoY + labelHeight + imgHeight + padding + 15
			addLabel(compositeImg, padding, formulaY, fmt.Sprintf("Gen: %d Score: %.2f", generation+1, best.Score))

			yOffset := formulaY + 15
			drawWrappedFormula := func(prefix, formula string) {
				content := wordwrap.NewContent(prefix + formula)
				wrapper := wordwrap.NewSimpleWrapper([]*wordwrap.Content{content}, basicfont.Face7x13)
				lines, _, err := wrapper.TextToRect(image.Rect(0, 0, canvasWidth-2*padding, canvasHeight))
				if err != nil {
					log.Printf("Error wrapping text: %v", err)
					return
				}
				for _, line := range lines {
					if yOffset > canvasHeight-20 {
						break
					}
					addLabel(compositeImg, padding, yOffset, line.TextValue())
					yOffset += 15
				}
			}

			drawWrappedFormula("R: ", best.Rf.String())
			drawWrappedFormula("G: ", best.Gf.String())
			drawWrappedFormula("B: ", best.Bf.String())

			// Convert to Paletted for GIF
			palettedImg := image.NewPaletted(compositeRect, palette.Plan9)
			draw.FloydSteinberg.Draw(palettedImg, compositeRect, compositeImg, image.Pt(0, 0))

			frames = append(frames, palettedImg)
			delays = append(delays, 20) // 200ms delay per frame
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

func addLabel(img *image.RGBA, x, y int, label string) {
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

