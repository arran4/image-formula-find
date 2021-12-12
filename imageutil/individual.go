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
	DNA    string
	Parent []*Individual
	Score  float64
	rf     *image_formula_find.Function
	bf     *image_formula_find.Function
	gf     *image_formula_find.Function
	i      draw.Image
	d      *drawer1.Drawer
}

func (i *Individual) Calculate(plotSize image.Rectangle, srcimg image.Image) {
	rd, bd, gd := dna1.SplitString3(i.DNA)
	i.rf = dna1.ParseFunction(rd)
	i.bf = dna1.ParseFunction(bd)
	i.gf = dna1.ParseFunction(gd)
	i.d = &drawer1.Drawer{
		RedFormula:   i.rf,
		BlueFormula:  i.bf,
		GreenFormula: i.gf,
	}
	i.i = image.NewRGBA(plotSize.Bounds())
	draw.Draw(i.i, plotSize, i.d, image.Pt(0, 0), draw.Src)
	i.Score = CalculateDistance(srcimg, i.i)
}

func (i *Individual) CsvRow() []string {
	return []string{
		i.DNA, i.rf.String(), i.bf.String(), i.gf.String(), fmt.Sprintf("%0.2f", i.Score),
	}
}

func (i *Individual) Image() image.Image {
	return i.i
}
