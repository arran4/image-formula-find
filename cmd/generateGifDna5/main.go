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
	"time"

	"math"

	"github.com/arran4/golang-wordwrap"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"image-formula-find/dna5"
)

func main() {
	var inputPath string
	var outputPath string
	var generations int
	var steps int

	flag.StringVar(&inputPath, "input", "flag_space.png", "Path to input image")
	flag.StringVar(&outputPath, "output", "evolution-dna5.gif", "Path to output GIF")
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
		dnaBarHeight  = 120 // Increased height for wrapping
		borderWidth   = 2
	)

	canvasWidth := padding*2 + imgWidth
	if canvasWidth < 400 { canvasWidth = 400 } // Ensure width for text

	canvasHeight := padding + labelHeight + imgHeight + padding + labelHeight + imgHeight + padding + dnaBarHeight + padding + formulaHeight + padding

	// Using BasicRequired implementation from dna5
	worker := &dna5.BasicRequired{
		R: plotSize,
		I: srcimg,
	}

	var lastGeneration []*dna5.Individual
	newDNA := make(chan string, 100)
	go func() {
		for {
			dna := dna5.RndStr(50)
			if !dna5.Valid(dna) {
				continue
			}
			newDNA <- dna
		}
	}()

	var frames []*image.Paletted
	var delays []int

	stepInterval := generations / steps
	if stepInterval < 1 {
		stepInterval = 1
	}

	log.Printf("Starting evolution. Generations: %d, Steps: %d, Interval: %d", generations, steps, stepInterval)

	startTime := time.Now()
	lastLogTime := time.Now()

	for generation := 0; generation < generations; generation++ {
		lastGeneration = dna5.GenerationProcess(worker, lastGeneration, generation, newDNA)

		// Log progress every 10 seconds to keep CI alive and inform user
		if time.Since(lastLogTime) > 10*time.Second {
			elapsed := time.Since(startTime)
			avgTimePerGen := elapsed / time.Duration(generation+1)
			remainingGens := generations - (generation + 1)
			estimatedRemaining := avgTimePerGen * time.Duration(remainingGens)

			log.Printf("Generation %d/%d (%.2f%%). Elapsed: %v. Estimated remaining: %v",
				generation+1, generations, float64(generation+1)/float64(generations)*100,
				elapsed.Round(time.Second), estimatedRemaining.Round(time.Second))
			lastLogTime = time.Now()
		}

		if (generation+1)%stepInterval == 0 || generation == generations-1 {
			log.Printf("Capturing frame at generation %d", generation+1)

			best := lastGeneration[0]
			evolvedImg := best.Image()

			compositeRect := image.Rect(0, 0, canvasWidth, canvasHeight)
			compositeImg := image.NewRGBA(compositeRect)

			draw.Draw(compositeImg, compositeRect, &image.Uniform{color.White}, image.Pt(0, 0), draw.Src)

			drawBorder := func(r image.Rectangle, c color.Color) {
				draw.Draw(compositeImg, image.Rect(r.Min.X-borderWidth, r.Min.Y-borderWidth, r.Max.X+borderWidth, r.Min.Y), &image.Uniform{c}, image.Pt(0, 0), draw.Src)
				draw.Draw(compositeImg, image.Rect(r.Min.X-borderWidth, r.Max.Y, r.Max.X+borderWidth, r.Max.Y+borderWidth), &image.Uniform{c}, image.Pt(0, 0), draw.Src)
				draw.Draw(compositeImg, image.Rect(r.Min.X-borderWidth, r.Min.Y, r.Min.X, r.Max.Y), &image.Uniform{c}, image.Pt(0, 0), draw.Src)
				draw.Draw(compositeImg, image.Rect(r.Max.X, r.Min.Y, r.Max.X+borderWidth, r.Max.Y), &image.Uniform{c}, image.Pt(0, 0), draw.Src)
			}

			imgX := (canvasWidth - imgWidth) / 2

			addLabel(compositeImg, padding, padding+labelHeight-5, "Target")
			targetRect := image.Rect(imgX, padding+labelHeight, imgX+imgWidth, padding+labelHeight+imgHeight)
			drawBorder(targetRect, color.Black)
			draw.Draw(compositeImg, targetRect, srcimg, image.Pt(0, 0), draw.Src)

			evoY := padding + labelHeight + imgHeight + padding
			addLabel(compositeImg, padding, evoY+labelHeight-5, "Evolution")
			evolvedRect := image.Rect(imgX, evoY+labelHeight, imgX+imgWidth, evoY+labelHeight+imgHeight)
			drawBorder(evolvedRect, color.Black)
			draw.Draw(compositeImg, evolvedRect, evolvedImg, image.Pt(0, 0), draw.Src)

			dnaY := evoY + labelHeight + imgHeight + padding
			dnaRect := image.Rect(padding, dnaY, canvasWidth-padding, dnaY+dnaBarHeight)
			drawDNABar(compositeImg, dnaRect, best.DNA)

			formulaY := dnaY + dnaBarHeight + padding + 15
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

			palettedImg := image.NewPaletted(compositeRect, palette.Plan9)
			draw.FloydSteinberg.Draw(palettedImg, compositeRect, compositeImg, image.Pt(0, 0))

			frames = append(frames, palettedImg)
			delays = append(delays, 20)
		}
	}

	fout, err := os.Create(outputPath)
	if err != nil {
		log.Panicf("Error: %v", err)
	}
	defer func() {
		if err := fout.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

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
	defer func() {
		if err := fin.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()
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

func drawDNABar(img *image.RGBA, r image.Rectangle, dna string) {
	// Split DNA into R, G, B channels using DNA5
	rStr, gStr, bStr := dna5.SplitString3(dna)

	drawChannelBar := func(rect image.Rectangle, s string, baseColor color.RGBA) {
		if len(s) == 0 {
			return
		}

		const blockWidth = 4
		const rowHeight = 10

		blocksPerRow := rect.Dx() / blockWidth
		if blocksPerRow < 1 { blocksPerRow = 1 }

		for i, char := range s {
			row := i / blocksPerRow
			col := i % blocksPerRow

			x1 := rect.Min.X + col*blockWidth
			y1 := rect.Min.Y + row*rowHeight
			x2 := x1 + blockWidth
			y2 := y1 + rowHeight

			if y2 > rect.Max.Y {
				break
			}

			hue := (int(char) * 10) % 360
			c := hsvToRGB(float64(hue), 1.0, 1.0)

			draw.Draw(img, image.Rect(x1, y1, x2, y2), &image.Uniform{c}, image.Pt(0, 0), draw.Src)
		}
	}

	h := r.Dy() / 3
	drawChannelBar(image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+h), rStr, color.RGBA{0, 0, 0, 255})
	drawChannelBar(image.Rect(r.Min.X, r.Min.Y+h, r.Max.X, r.Min.Y+2*h), gStr, color.RGBA{0, 0, 0, 255})
	drawChannelBar(image.Rect(r.Min.X, r.Min.Y+2*h, r.Max.X, r.Max.Y), bStr, color.RGBA{0, 0, 0, 255})
}

func hsvToRGB(h, s, v float64) color.RGBA {
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := v - c
	var r, g, b float64
	if h >= 0 && h < 60 {
		r, g, b = c, x, 0
	} else if h >= 60 && h < 120 {
		r, g, b = x, c, 0
	} else if h >= 120 && h < 180 {
		r, g, b = 0, c, x
	} else if h >= 180 && h < 240 {
		r, g, b = 0, x, c
	} else if h >= 240 && h < 300 {
		r, g, b = x, 0, c
	} else {
		r, g, b = c, 0, x
	}
	return color.RGBA{
		R: uint8((r + m) * 255),
		G: uint8((g + m) * 255),
		B: uint8((b + m) * 255),
		A: 255,
	}
}
