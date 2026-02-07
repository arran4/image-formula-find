package image_formula_find

import (
	"testing"
)

func BenchmarkEvaluate(b *testing.B) {
	exprStr := "y / 4 = x * x + 2"
	f, err := ParseFunction(exprStr)
	if err != nil {
		b.Fatal(err)
	}
	if f == nil {
		b.Fatal("Failed to parse function")
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = f.Evaluate(float64(i), float64(i), i)
	}
}
