//go:build noebiten

package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"time"
)

type EbitenImage struct{
	Img image.Image
}
type EbitenDrawImageOptions struct{}
type EbitenGame interface {
	Update() error
	Draw(screen *EbitenImage)
	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
}

func EbitenSetWindowSize(width, height int) {}

func EbitenSetWindowTitle(title string) {}

func EbitenRunGame(game EbitenGame) error {
	// Headless loop logic
	// In the original ebiten.RunGame, it loops and calls Update and Draw.
	// Here we will implement a ticker loop similar to what I did in `noebiten.go` before.

	g, ok := game.(*Game)
	if !ok {
		return fmt.Errorf("game is not of type *Game")
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Println("Running in headless mode (mocking Ebiten)")

	// Create a dummy screen for Draw calls
	dummyScreen := &EbitenImage{}

	for {
		<-ticker.C

		w := g.Worker
		w.RLock()
		gen := w.Generation
		var bestImg image.Image
		if len(w.Winners) > 0 {
			bestImg = w.Winners[0].Image()
		} else if len(w.LastGeneration) > 0 {
			bestImg = w.LastGeneration[0].Image()
		}
		w.RUnlock()

		if bestImg != nil {
			filename := fmt.Sprintf("gen_%d.png", gen)
			f, err := os.Create(filename)
			if err != nil {
				log.Printf("Failed to create file: %v", err)
			} else {
				if err := png.Encode(f, bestImg); err != nil {
					log.Printf("Failed to encode png: %v", err)
				}
				f.Close()
				log.Printf("Saved %s", filename)
			}
		}

		// Simulate game loop calls
		if err := game.Update(); err != nil {
			return err
		}
		game.Draw(dummyScreen)
	}
}

func EbitenNewImageFromImage(source image.Image) *EbitenImage {
	return &EbitenImage{Img: source}
}

func EbitenDebugPrint(image *EbitenImage, str string) {}

func EbitenDebugPrintAt(image *EbitenImage, str string, x, y int) {}

func EbitenTranslate(op *EbitenDrawImageOptions, x, y float64) {}

func EbitenNewDrawImageOptions() *EbitenDrawImageOptions {
	return &EbitenDrawImageOptions{}
}

// Add method to EbitenImage to satisfy Draw call `screen.DrawImage`
func (i *EbitenImage) DrawImage(img *EbitenImage, op *EbitenDrawImageOptions) {
	// No-op
}

func EbitenDrawImage(dst *EbitenImage, src *EbitenImage, op *EbitenDrawImageOptions) {
	dst.DrawImage(src, op)
}
