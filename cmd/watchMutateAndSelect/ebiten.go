//go:build !noebiten

package main

import (
	"fmt"
	"image-formula-find/worker"
	"log"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func run(w *worker.Worker) {
	game := NewGame(w)
	ebiten.SetWindowSize(640*2, 480*3)
	ebiten.SetWindowTitle("Watch Generator")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	*worker.Worker
	esrcimg *ebiten.Image
}

func (game *Game) Update() error {
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	game.RLock()
	defer game.RUnlock()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Generation %d", game.Generation))
	op := &ebiten.DrawImageOptions{}
	const offsetY = 20
	op.GeoM.Translate(0, offsetY)
	if game.esrcimg == nil {
		game.esrcimg = ebiten.NewImageFromImage(game.SourceImage())
	}
	screen.DrawImage(game.esrcimg, op)
	for i, ind := range game.LastGeneration {
		op := &ebiten.DrawImageOptions{}
		tx := 0 //(worker.srcimg.Bounds().Dx()+10)
		ty := (game.SourceImage().Bounds().Dy() + 10) * (i + 1)
		op.GeoM.Translate(float64(tx), float64(ty))
		i := ind.Image()
		if i != nil {
			screen.DrawImage(ebiten.NewImageFromImage(i), op)
		}
		tx += (game.SourceImage().Bounds().Dx() + 10)
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

func (game *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func NewGame(w *worker.Worker) *Game {
	return &Game{
		Worker: w,
	}
}
