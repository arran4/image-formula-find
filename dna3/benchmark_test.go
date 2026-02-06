package dna3

import (
	"fmt"
	"testing"
)

func BenchmarkRndStr(b *testing.B) {
	lengths := []int{10, 100, 1000, 10000}

	for _, length := range lengths {
		b.Run(fmt.Sprintf("Length-%d", length), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = RndStr(length)
			}
		})
	}
}
