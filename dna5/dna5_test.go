package dna5

import (
	"image"
	"image-formula-find"
	"testing"
)

type TestRequired struct {
	R image.Rectangle
	I image.Image
}

func (b *TestRequired) PlotSize() image.Rectangle {
	return b.R
}

func (b *TestRequired) SourceImage() image.Image {
	return b.I
}

func TestSplitString3(t *testing.T) {
	// 0 1 2 0 1 2
	// R G B R G B
	s := "ABCDEF"
	// R: A, D
	// G: B, E
	// B: C, F
	r, g, b := SplitString3(s)
	if r != "AD" {
		t.Errorf("Expected AD, got %s", r)
	}
	if g != "BE" {
		t.Errorf("Expected BE, got %s", g)
	}
	if b != "CF" {
		t.Errorf("Expected CF, got %s", b)
	}

	// Test Uneven
	s = "ABCDE"
	// R: A, D
	// G: B, E
	// B: C
	r, g, b = SplitString3(s)
	if r != "AD" {
		t.Errorf("Expected AD, got %s", r)
	}
	if g != "BE" {
		t.Errorf("Expected BE, got %s", g)
	}
	if b != "C" {
		t.Errorf("Expected C, got %s", b)
	}
}

func TestParseRPN(t *testing.T) {
	// A=0(X), B=1(Y), Q=16(+)
	// DNA: "ABQ" -> X Y +
	dna := "ABQ"
	expr := ParseRPN(dna)
	if expr.String() != "X + Y" {
		t.Errorf("Expected X + Y, got %s", expr.String())
	}

	// Constants
	// D=3(0.1), E=4(-0.1)
	// DNA: "DE" -> 0.1 -0.1 (Top of stack is -0.1)
	dna = "DE"
	expr = ParseRPN(dna)
	// Should be -0.1
	if expr.String() != "-0.1" {
		t.Errorf("Expected -0.1, got %s", expr.String())
	}

	// Complex: A B Q C R -> (X+Y)*T
	// A=X, B=Y, Q=+, C=T, R=17(-) (Wait, 18 is *)
	// S=18
	// DNA: "ABQS" -> X Y + * (Stack: X+Y. Need one more op. * needs 2 args.
	// Stack: [X+Y]. * sees 1 arg, pop X+Y. Pushes X+Y back?
	// Code: rhs=pop, lhs=pop. if lhs!=nil && rhs!=nil { push(lhs*rhs) } else if rhs!=nil { push(rhs) }
	// Stack has [X+Y]. rhs = X+Y. lhs = nil. Push(rhs) -> [X+Y].
	// So * is no-op if stack has 1 item.

	// Let's try: A B Q C S -> X Y + T * -> (X+Y)*T
	// A=0, B=1, Q=16, C=2(T), S=18(*)
	dna = "ABQCS"
	expr = ParseRPN(dna)
	// Note: String() does not add parentheses for precedence, so (X+Y)*T renders as "X + Y * T"
	// We verify the structure to be sure.
	if _, ok := expr.(*image_formula_find.Multiply); !ok {
		t.Errorf("Expected Multiply root, got %T", expr)
	}
	if expr.String() != "X + Y * T" {
		t.Errorf("Expected X + Y * T, got %s", expr.String())
	}

	// Test New Functions
	// j=35 (Sinh), k=36 (Cosh)
	// DNA: "Aj" -> Sinh(X)
	dna = "Aj"
	expr = ParseRPN(dna)
	if expr.String() != "Sinh(X)" {
		t.Errorf("Expected Sinh(X), got %s", expr.String())
	}

	// Atan2 (41): p(41)
	// DNA: "ABp" -> Atan2(X, Y)
	dna = "ABp"
	expr = ParseRPN(dna)
	if expr.String() != "Atan2(X, Y)" {
		t.Errorf("Expected Atan2(X, Y), got %s", expr.String())
	}
}

func TestGenerationProcess(t *testing.T) {
	// Mock Required
	req := &BasicRequired{
		R: image.Rect(0, 0, 10, 10),
		I: image.NewRGBA(image.Rect(0, 0, 10, 10)),
	}

	newDNA := make(chan string, 100)
	go func() {
		// Provide some valid DNA to ensure we don't block
		// Needs X(A) and Y(B) in all 3 splits.
		// A long string of "ABQ" (X Y +) repeated should work and split into chunks having X and Y.
		validDNA := ""
		for k := 0; k < 20; k++ {
			validDNA += "ABQ"
		}
	    for i := 0; i < 100; i++ {
	        newDNA <- validDNA // Ensure valid DNA is available
			newDNA <- RndStr(50)
	    }
		close(newDNA)
	}()

	gen := GenerationProcess(req, nil, 0, newDNA)
	if len(gen) == 0 {
		t.Error("Generation produced no children")
	}

	gen2 := GenerationProcess(req, gen, 1, newDNA)
	if len(gen2) == 0 {
		t.Error("Gen 2 produced no children")
	}
}
