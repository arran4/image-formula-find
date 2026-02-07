package image_formula_find

import (
	"testing"
)

func BenchmarkSingleFunctionEvaluate(b *testing.B) {
	// Setup
	sf := NewSingleFunction("Sin", &Var{Var: "x"})
	state := &State{X: 0.5, Y: 0.5, T: 1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sf.Evaluate(state)
	}
}

func BenchmarkDoubleFunctionEvaluate(b *testing.B) {
	// Setup
	df := NewDoubleFunction("Pow", &Var{Var: "x"}, &Const{Value: 2}, false)
	state := &State{X: 0.5, Y: 0.5, T: 1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		df.Evaluate(state)
	}
}
