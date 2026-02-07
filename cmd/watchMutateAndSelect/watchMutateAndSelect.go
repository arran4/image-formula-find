package main

import (
	"fmt"
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
	game := NewGame()
	go game.Work()
	EbitenSetWindowSize(640*2, 480*3)
	EbitenSetWindowTitle("Watch Generator")
	if err := EbitenRunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	*worker.Worker
	esrcimg *EbitenImage
}

func (game *Game) Update() error {
	return nil
}

func (game *Game) Draw(screen *EbitenImage) {
	game.RLock()
	defer game.RUnlock()
	EbitenDebugPrint(screen, fmt.Sprintf("Generation %d", game.Generation))
	op := EbitenNewDrawImageOptions()
	const offsetY = 20
	EbitenTranslate(op, 0, offsetY)
	if game.esrcimg == nil {
		game.esrcimg = EbitenNewImageFromImage(game.SourceImage())
	}
	screen.DrawImage(game.esrcimg, op)
	for i, ind := range game.LastGeneration {
		op := EbitenNewDrawImageOptions()
		tx := 0 //(worker.srcimg.Bounds().Dx()+10)
		ty := (game.SourceImage().Bounds().Dy() + 10) * (i + 1)
		EbitenTranslate(op, float64(tx), float64(ty))
		i := ind.Image()
		if i != nil {
			screen.DrawImage(EbitenNewImageFromImage(i), op)
		}
		tx += (game.SourceImage().Bounds().Dx() + 10)
		EbitenDebugPrintAt(screen, fmt.Sprintf("Score %f", ind.Score), tx, ty)
		ty += 20
		EbitenDebugPrintAt(screen, fmt.Sprintf("DNA %s", ind.DNA), tx, ty)
		ty += 20
		EbitenDebugPrintAt(screen, fmt.Sprintf("Red   = %s", ind.Rf.String()), tx, ty)
		ty += 20
		EbitenDebugPrintAt(screen, fmt.Sprintf("Green = %s", ind.Gf.String()), tx, ty)
		ty += 20
		EbitenDebugPrintAt(screen, fmt.Sprintf("Blue  = %s", ind.Bf.String()), tx, ty)
	}
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func NewGame() *Game {
	w := worker.NewWorker(LoadImage())
	return &Game{
		Worker: w,
	}
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
