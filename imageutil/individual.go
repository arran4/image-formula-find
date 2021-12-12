package imageutil

import (
	"fmt"
	"image"
	"image-formula-find"
	"image-formula-find/dna1"
	"image-formula-find/drawer1"
	"image/draw"
)

type Sorter struct {
	Children []*Individual
}

func (s *Sorter) Len() int {
	return len(s.Children)
}

func (s *Sorter) Less(i, j int) bool {
	return s.Children[i].Score < s.Children[j].Score
}

func (s *Sorter) Swap(i, j int) {
	s.Children[i], s.Children[j] = s.Children[j], s.Children[i]
}

type Individual struct {
	DNA             string
	Parent          []*Individual
	Score           float64
	Rf              *image_formula_find.Function
	Bf              *image_formula_find.Function
	Gf              *image_formula_find.Function
	i               draw.Image
	d               *drawer1.Drawer
	FirstGeneration int
}

type Required interface {
	PlotSize() image.Rectangle
	SourceImage() image.Image
}

type BasicRequired struct {
	R image.Rectangle
	I image.Image
}

func (b *BasicRequired) PlotSize() image.Rectangle {
	return b.R
}

func (b *BasicRequired) SourceImage() image.Image {
	return b.I
}

func (i *Individual) Calculate(required Required) {
	i.Rf, i.Bf, i.Gf = dna1.ParseDNA(i.DNA)
	i.d = &drawer1.Drawer{
		RedFormula:   i.Rf,
		BlueFormula:  i.Bf,
		GreenFormula: i.Gf,
	}
	i.i = image.NewRGBA(required.PlotSize().Bounds())
	draw.Draw(i.i, required.PlotSize(), i.d, image.Pt(0, 0), draw.Src)
	i.Score = CalculateDistance(required.SourceImage(), i.i)
}

func (i *Individual) CsvRow() []string {
	return []string{
		i.DNA, i.Rf.String(), i.Bf.String(), i.Gf.String(), fmt.Sprintf("%0.2f", i.Score),
	}
}

func (i *Individual) Image() image.Image {
	return i.i
}
