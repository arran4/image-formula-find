package dna1

import (
	"math/rand"
	"testing"
)

// breedOriginal is the original implementation of Breed for benchmarking comparison.
func breedOriginal(a string, b string) string {
	p := 10
	if len(a) < p {
		p = len(a)
	}
	if len(b) < p {
		p = len(b)
	}
	result := ""
	for i := 0; i < p; i++ {
		var s string
		switch rand.Int31n(2) {
		case 0:
			s = a
		case 1:
			s = b
		}
		st := (len(s) / p) * i
		e := (len(s)/p)*(i+1) - 1
		result += s[st:e]
	}
	return result
}

func BenchmarkBreedOriginal(b *testing.B) {
	// Setup inputs
	s1 := RndStr(100)
	s2 := RndStr(100)

	// Reset timer to ignore setup time
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		breedOriginal(s1, s2)
	}
}

func BenchmarkBreed(b *testing.B) {
	// Setup inputs
	s1 := RndStr(100)
	s2 := RndStr(100)

	// Reset timer to ignore setup time
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Breed(s1, s2)
	}
}

func BenchmarkBreedShort(b *testing.B) {
	// Setup inputs
	s1 := RndStr(10)
	s2 := RndStr(10)

	// Reset timer to ignore setup time
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Breed(s1, s2)
	}
}

func BenchmarkBreedLong(b *testing.B) {
	// Setup inputs
	s1 := RndStr(1000)
	s2 := RndStr(1000)

	// Reset timer to ignore setup time
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Breed(s1, s2)
	}
}

func TestBreedProperties(t *testing.T) {
	// Use length that is multiple of 10 to make math easy?
	// p=10. len=100. len/p = 10.
	// Chunk size = 10 * (i+1) - 1 - (10*i) = 10*i + 10 - 1 - 10*i = 9.
	// Wait: st = 10*i. e = 10*(i+1) - 1.
	// Slice is s[st:e]. Length is e - st = 10(i+1) - 1 - 10i = 9.
	// So expected total length is 9 * 10 = 90.

	length := 100
	s1 := RndStr(length)
	s2 := RndStr(length)

	// Ensure they are different enough to be useful
	for s1 == s2 {
		s2 = RndStr(length)
	}

	res := Breed(s1, s2)

	// Calculate expected chunk length
	p := 10
	chunkLen := (length/p) - 1
	expectedTotalLen := chunkLen * p

	if len(res) != expectedTotalLen {
		t.Errorf("Expected length %d, got %d", expectedTotalLen, len(res))
	}

	// Verify each chunk
	for i := 0; i < p; i++ {
		start := i * chunkLen
		end := (i + 1) * chunkLen
		if end > len(res) {
			end = len(res)
		}

		chunk := res[start:end]

		// Reconstruct expected chunk from s1
		st1 := (len(s1) / p) * i
		e1 := (len(s1)/p)*(i+1) - 1
		chunk1 := s1[st1:e1]

		// Reconstruct expected chunk from s2
		st2 := (len(s2) / p) * i
		e2 := (len(s2)/p)*(i+1) - 1
		chunk2 := s2[st2:e2]

		if chunk != chunk1 && chunk != chunk2 {
			t.Errorf("Chunk %d (%q) does not match s1 chunk (%q) or s2 chunk (%q)", i, chunk, chunk1, chunk2)
		}
	}
}
