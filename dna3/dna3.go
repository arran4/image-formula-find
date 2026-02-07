package dna3

import (
	"image-formula-find"
	"math"
	"math/rand"
	"sort"
	"strings"
	"sync"

	"github.com/agnivade/levenshtein"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

var (
	runeMapPos = map[rune]int{}
)

func init() {
	for p, c := range chars {
		runeMapPos[c] = p
	}
}

func RndStr(length int) string {
	var sb strings.Builder
	sb.Grow(length)
	for i := 0; i < length; i++ {
		sb.WriteByte(chars[rand.Int31n(int32(len(chars)))])
	}
	return sb.String()
}

// ParseDNA splits the DNA into 3 channels (R, G, B) and parses each using the dna3 logic.
func ParseDNA(dna string) (*image_formula_find.Function, *image_formula_find.Function, *image_formula_find.Function) {
	rd, bd, gd := SplitString3(dna)
	rf := ParseFunction(rd)
	bf := ParseFunction(bd)
	gf := ParseFunction(gd)
	return rf, bf, gf
}

func ParseFunction(arg string) *image_formula_find.Function {
	expr := ParseChannel(arg)
	return &image_formula_find.Function{
		Equals: &image_formula_find.Equals{
			LHS: &image_formula_find.Var{Var: "Y"},
			RHS: expr,
		},
	}
}

func isZero(expr image_formula_find.Expression) bool {
	if c, ok := expr.(*image_formula_find.Const); ok {
		return c.Value == 0
	}
	return false
}

// ParseChannel implements the fixed structure formula builder.
// Formula = Sum of 6 Terms.
// Term k (0..5) uses params P_{5k}..P_{5k+4}.
// If k is even: Term = X * P_4^P_3 + P_2*P_1 + P_0
// If k is odd:  Term = Y * P_4^P_3 + P_2*P_1 + P_0
func ParseChannel(dna string) image_formula_find.Expression {
	var totalExpr image_formula_find.Expression = &image_formula_find.Const{Value: 0}

	for k := 0; k < 6; k++ {
		baseIdx := k * 5

		p0 := Resolve(baseIdx, dna)
		p1 := Resolve(baseIdx+1, dna)
		p2 := Resolve(baseIdx+2, dna)
		p3 := Resolve(baseIdx+3, dna)
		p4 := Resolve(baseIdx+4, dna)

		// Term Structure:
		// Part A: P_2 * P_1
		var partA image_formula_find.Expression
		if isZero(p2) || isZero(p1) {
			partA = &image_formula_find.Const{Value: 0}
		} else {
			partA = &image_formula_find.Multiply{
				LHS: p2,
				RHS: p1,
			}
		}

		// Part B: X/Y * P_4^P_3
		var partB image_formula_find.Expression
		if isZero(p4) {
			partB = &image_formula_find.Const{Value: 0}
		} else {
			var variable image_formula_find.Expression
			if k%2 == 0 {
				variable = &image_formula_find.Var{Var: "X"}
			} else {
				variable = &image_formula_find.Var{Var: "Y"}
			}

			powerPart := &image_formula_find.Power{
				LHS: p4,
				RHS: p3,
			}

			partB = &image_formula_find.Multiply{
				LHS: variable,
				RHS: powerPart,
			}
		}

		term := p0
		if isZero(p0) {
			term = &image_formula_find.Const{Value: 0}
		}

		if !isZero(partA) {
			if isZero(term) {
				term = partA
			} else {
				term = &image_formula_find.Plus{
					LHS: partA,
					RHS: term,
				}
			}
		}

		if !isZero(partB) {
			if isZero(term) {
				term = partB
			} else {
				term = &image_formula_find.Plus{
					LHS: partB,
					RHS: term,
				}
			}
		}

		if !isZero(term) {
			if isZero(totalExpr) {
				totalExpr = term
			} else {
				totalExpr = &image_formula_find.Plus{
					LHS: totalExpr,
					RHS: term,
				}
			}
		}
	}
	return totalExpr
}

// Resolve returns the expression for parameter at index `idx` considering layers.
func Resolve(idx int, dna string) image_formula_find.Expression {
	if idx >= len(dna) {
		return &image_formula_find.Const{Value: 0}
	}

	// Layer 0: Base Value
	char := dna[idx]
	val := MapValue(char)
	var expr image_formula_find.Expression = &image_formula_find.Const{Value: val}

	// Apply layers
	const Period = 30
	layer := 1
	for {
		nextIdx := idx + layer*Period
		if nextIdx >= len(dna) {
			break
		}

		c := dna[nextIdx]
		v := runeMapPos[rune(c)]

		if v == 0 {
			// 'A' -> Empty/Skip
			layer++
			continue
		}

		if v >= 54 {
			// Op
			opName := MapOp(v)
			expr = &image_formula_find.SingleFunction{
				Name: opName,
				Expr: expr,
			}
		} else {
			// Number -> Add
			// Note: We use MapValue again here, or should we use a simpler addition?
			// Using MapValue keeps consistency.
			numVal := MapValue(c)
			if numVal != 0 {
				expr = &image_formula_find.Plus{
					LHS: expr,
					RHS: &image_formula_find.Const{Value: numVal},
				}
			}
		}
		layer++
	}

	return expr
}

func MapValue(c uint8) float64 {
	i := runeMapPos[rune(c)]
	// 0 (A) is Empty/Zero
	if i == 0 {
		return 0
	}

	// Range 1..53 (B..1)
	// Requested: "range from 0 to 500 by step by 25"
	// Indices 1..21 map to positive 25, 50 ... 500
	// Indices 22..42 map to negative -25, -50 ... -500
	// Indices 43..53 -> Empty/Zero to reduce constants

	if i <= 21 {
		return float64(i) * 25.0
	} else if i <= 42 {
		return -float64(i-21) * 25.0
	} else {
		return 0
	}
}

func MapOp(v int) string {
	switch v {
	case 54: return "Sin"
	case 55: return "Cos"
	case 56: return "Tan"  // Warning: Asymptotes
	case 57: return "Atan"
	case 58: return "Abs"
	case 59: return "Exp" // Warning: Growth
	case 60: return "Log" // Warning: Negative domain
	case 61: return "Sqrt" // Warning: Negative domain
	case 62: return "Negate"
	case 63: return "Sin" // Extra slot?
	default: return "Sin"
	}
}

// Helpers reused from dna1

func SplitString3(arg string) (string, string, string) {
	p1 := ""
	p2 := ""
	p3 := ""
	for len(arg) > 1 {
		c1 := arg[0]
		arg = arg[1:]
		i1, ok := runeMapPos[rune(c1)]
		if !ok {
			continue
		}
		var i2 int
		for i2-i1 < 2 && len(arg) > 0 {
			c2 := arg[0]
			arg = arg[1:]
			i2, ok = runeMapPos[rune(c2)]
			if !ok {
				continue
			}
			if i2 < i1 {
				i1, i2 = i2, i1
			}
		}
		width := float64(len(arg) - 3)
		if width < 0 {
			break
		}
		incs := width / 64.0
		if int(math.Round(incs*float64(i2))) > len(arg) {
			break
		}
		p1 = arg[0 : int(math.Round(incs*float64(i1)))+1]
		p2 = arg[int(math.Round(incs*float64(i1)))+1 : 1+int(math.Round(incs*float64(i2)))+1]
		p3 = arg[1+int(math.Round(incs*float64(i2)))+1:]
		return p1, p2, p3
	}
	return p1, p2, p3
}

func Mutate(a string) string {
	switch rand.Int31n(12) {
	case 0:
		return AppendMutate(a)
	case 1:
		return PopMutate(a)
	case 2:
		return ShiftMutate(a)
	case 3:
		return UnshiftMutate(a)
	case 4:
		return DeleteMutate(a)
	case 5:
		return InsertMutate(a)
	default:
		return PositionMutate(a)
	}
}

func UnshiftMutate(a string) string {
	return string([]byte{chars[rand.Int31n(int32(len(chars)))]}) + a
}

func InsertMutate(a string) string {
	if len(a) == 0 {
		return a
	}
	p := int(rand.Int31n(int32(len(a))))
	return a[:p] + string([]byte{chars[rand.Int31n(int32(len(chars)))]}) + a[p:]
}

func DeleteMutate(a string) string {
	if len(a) == 0 {
		return a
	}
	p := int(rand.Int31n(int32(len(a))))
	return a[:p] + a[p+1:]
}

func ShiftMutate(a string) string {
	if len(a) == 0 {
		return a
	}
	return a[1:]
}

func PopMutate(a string) string {
	if len(a) == 0 {
		return a
	}
	return a[:len(a)-1]
}

func AppendMutate(a string) string {
	return a + string([]byte{chars[rand.Int31n(int32(len(chars)))]})
}

func PositionMutate(a string) string {
	if len(a) == 0 {
		return a
	}
	p := int(rand.Int31n(int32(len(a))))
	return a[:p] + string([]byte{chars[rand.Int31n(int32(len(chars)))]}) + a[p+1:]
}

func Breed(a string, b string) string {
	p := 10
	if len(a) < p {
		p = len(a)
	}
	if len(b) < p {
		p = len(b)
	}
	var result strings.Builder
	// Pre-allocate to avoid reallocations.
	// We use the maximum length of a and b as a safe upper bound.
	growLen := len(a)
	if len(b) > growLen {
		growLen = len(b)
	}
	result.Grow(growLen)

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
		result.WriteString(s[st:e])
	}
	return result.String()
}

func Valid(dna string) bool {
    // dna3 produces formulas that always have X and Y if length is sufficient.
    // Even if length is 0, ParseChannel returns Const(0).
    // So technically it's always "valid" structure-wise.
    // But image-formula-find expects non-nil.
    // ParseDNA handles it.
	return true
}

const (
	childrenCount = 10
)

func GenerationProcess(worker Required, lastGeneration []*Individual, generation int, newDNA chan string) []*Individual {
	const mutations = 8
	var children = make([]*Individual, 0, len(lastGeneration)*mutations+len(lastGeneration)*len(lastGeneration)+childrenCount+1)

	seen := map[string]struct{}{}

	for _, p := range lastGeneration {
		dna := p.DNA
		if _, ok := seen[dna]; ok {
			continue
		}
		seen[dna] = struct{}{}
		children = append(children, p)
	}

	for _, p := range lastGeneration {
		for i := 0; i <= mutations; i++ {
			dna := p.DNA
			for m := 0; m <= i; m++ {
				dna = Mutate(dna)
			}
			if _, ok := seen[dna]; ok {
				continue
			}
			seen[dna] = struct{}{}
			if !Valid(dna) {
				continue
			}
			children = append(children, &Individual{
				DNA:             dna,
				Parent:          []*Individual{p},
				FirstGeneration: generation,
				Lineage:         p.DNA,
			})
		}
	}

	if len(lastGeneration) > 4 {
		p1 := lastGeneration[int(rand.Int31n(int32(len(lastGeneration))))]
		p2 := lastGeneration[int(rand.Int31n(int32(len(lastGeneration))))]
		dna := Breed(p1.DNA, p2.DNA)
		if _, ok := seen[dna]; !ok {
			seen[dna] = struct{}{}

			if Valid(dna) {
				if dna != p1.DNA && dna != p2.DNA && p1.Lineage != p2.Lineage {
					children = append(children, &Individual{
						DNA: dna,
						Parent: []*Individual{
							p1, p2,
						},
						Lineage:         dna,
						FirstGeneration: generation,
					})
				}
			}
		}
	}

	for len(children) < childrenCount {
		dna := <-newDNA
		if _, ok := seen[dna]; ok {
			continue
		}
		seen[dna] = struct{}{}
		if !Valid(dna) {
			continue
		}
		children = append(children, &Individual{
			DNA:             dna,
			Lineage:         dna,
			FirstGeneration: generation,
		})
	}

	wg := sync.WaitGroup{}
	for fi := range children {
		wg.Add(1)
		go func(i int, child *Individual) {
			defer wg.Done()
			child.Calculate(worker)
		}(fi, children[fi])
	}
	wg.Wait()

	sort.Sort((&Sorter{
		Children: children,
	}))

	lastGeneration = make([]*Individual, 0, childrenCount)
	for len(lastGeneration) < childrenCount && len(children) > 0 {
		child := children[0]
		children = children[1:]

		minDistance := math.MaxInt
		for _, lg := range lastGeneration {
			m := levenshtein.ComputeDistance(lg.DNA, child.DNA)
			if m < minDistance {
				minDistance = m
				if minDistance < 10 {
					break
				}
			}
		}

		if minDistance < 10 {
			continue
		}

		lastGeneration = append(lastGeneration, child)
	}
	sort.Sort((&Sorter{
		Children: lastGeneration,
	}))
	return lastGeneration
}
