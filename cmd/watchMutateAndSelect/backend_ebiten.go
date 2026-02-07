//go:build !noebiten

package main

import (
	"image"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type EbitenImage = ebiten.Image
type EbitenDrawImageOptions = ebiten.DrawImageOptions
type EbitenGame = ebiten.Game

func EbitenSetWindowSize(width, height int) {
	ebiten.SetWindowSize(width, height)
}

func EbitenSetWindowTitle(title string) {
	ebiten.SetWindowTitle(title)
}

func EbitenRunGame(game EbitenGame) error {
	return ebiten.RunGame(game)
}

func EbitenNewImageFromImage(source image.Image) *EbitenImage {
	return ebiten.NewImageFromImage(source)
}

func EbitenDebugPrint(image *EbitenImage, str string) {
	ebitenutil.DebugPrint(image, str)
}

func EbitenDebugPrintAt(image *EbitenImage, str string, x, y int) {
	ebitenutil.DebugPrintAt(image, str, x, y)
}

func EbitenDrawImage(dst *EbitenImage, src *EbitenImage, op *EbitenDrawImageOptions) {
	dst.DrawImage(src, op)
}

func EbitenTranslate(op *EbitenDrawImageOptions, x, y float64) {
	op.GeoM.Translate(x, y)
}

func EbitenNewDrawImageOptions() *EbitenDrawImageOptions {
	return &ebiten.DrawImageOptions{}
}
