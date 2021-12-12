package imageutil

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"io"
	"log"
	"math"
	"testing"
)

var (
	//go:embed "in.png"
	inimg1bytes []byte
	inimg1      image.Image = MustReadImage(bytes.NewReader(inimg1bytes))
	//go:embed "in2.png"
	inimg2bytes []byte
	inimg2      image.Image = MustReadImage(bytes.NewReader(inimg2bytes))
)

func MustReadImage(r io.Reader) image.Image {
	i, _, err := image.Decode(r)
	if err != nil {
		log.Panicln(err)
	}
	return i
}

func TestCalculateDistance(t *testing.T) {
	tests := []struct {
		name string
		i1   image.Image
		i2   image.Image
		want float64
	}{
		{name: "Same iamge is same image", i1: inimg1, i2: inimg1, want: 0.0},
		{name: "Similar iamge is not the same image", i1: inimg1, i2: inimg2, want: 4.2949579029e+09},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateDistance(tt.i1, tt.i2)
			gotRounded := math.Round(got*100) / 100
			if gotRounded != math.Round(tt.want*100)/100 {
				t.Errorf("CalculateDistance() = %v, want %v", gotRounded, tt.want)
			}
		})
	}
}
