package dna5

import "testing"

func BenchmarkRndStr(b *testing.B) {
	// Length of DNA string typically used
	length := 100
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RndStr(length)
	}
}

func BenchmarkRndStrLarge(b *testing.B) {
	// Larger length to exaggerate the performance difference
	length := 10000
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RndStr(length)
	}
}
