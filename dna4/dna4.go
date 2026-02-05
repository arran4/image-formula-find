package dna4

import (
	"image-formula-find"
	"math"
	"math/rand"
	"sort"
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
	result := ""
	for len(result) < length {
		result += string([]byte{chars[rand.Int31n(int32(len(chars)))]})
	}
	return result
}

// ParseDNA splits the DNA into 3 channels (R, G, B) and parses each using the Stack Machine (RPN).
func ParseDNA(dna string) (*image_formula_find.Function, *image_formula_find.Function, *image_formula_find.Function) {
	rd, bd, gd := SplitString3(dna)
	rf := ParseFunction(rd)
	bf := ParseFunction(bd)
	gf := ParseFunction(gd)
	return rf, bf, gf
}

func ParseFunction(arg string) *image_formula_find.Function {
	expr := ParseRPN(arg)
	return &image_formula_find.Function{
		Equals: &image_formula_find.Equals{
			LHS: &image_formula_find.Var{Var: "Y"},
			RHS: expr,
		},
	}
}

func ParseRPN(arg string) image_formula_find.Expression {
	stack := []image_formula_find.Expression{}

	push := func(e image_formula_find.Expression) {
		stack = append(stack, e)
	}
	pop := func() image_formula_find.Expression {
		if len(stack) == 0 {
			return nil
		}
		e := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		return e
	}

	for _, char := range arg {
		idx, ok := runeMapPos[char]
		if !ok {
			continue
		}

		switch idx {
		// Vars (0-2)
		case 0:
			push(&image_formula_find.Var{Var: "X"})
		case 1:
			push(&image_formula_find.Var{Var: "Y"})
		case 2:
			push(&image_formula_find.Var{Var: "T"})

		// Consts (3-10)
		case 3:
			push(&image_formula_find.Const{Value: 0.1})
		case 4:
			push(&image_formula_find.Const{Value: -0.1})
		case 5:
			push(&image_formula_find.Const{Value: 1})
		case 6:
			push(&image_formula_find.Const{Value: -1})
		case 7:
			push(&image_formula_find.Const{Value: 0})
		case 8:
			push(&image_formula_find.Const{Value: 0.5})
		case 9:
			push(&image_formula_find.Const{Value: -0.5})
		case 10:
			push(&image_formula_find.Const{Value: math.Pi})

		// Binary Ops (16-21)
		case 16: // +
			rhs := pop()
			lhs := pop()
			if lhs != nil && rhs != nil {
				push(&image_formula_find.Plus{LHS: lhs, RHS: rhs})
			} else if rhs != nil {
				push(rhs)
			}
		case 17: // -
			rhs := pop()
			lhs := pop()
			if lhs != nil && rhs != nil {
				push(&image_formula_find.Subtract{LHS: lhs, RHS: rhs})
			} else if rhs != nil {
				push(rhs)
			}
		case 18: // *
			rhs := pop()
			lhs := pop()
			if lhs != nil && rhs != nil {
				push(&image_formula_find.Multiply{LHS: lhs, RHS: rhs})
			} else if rhs != nil {
				push(rhs)
			}
		case 19: // /
			rhs := pop()
			lhs := pop()
			if lhs != nil && rhs != nil {
				push(&image_formula_find.Divide{LHS: lhs, RHS: rhs})
			} else if rhs != nil {
				push(rhs)
			}
		case 20: // ^
			rhs := pop()
			lhs := pop()
			if lhs != nil && rhs != nil {
				push(&image_formula_find.Power{LHS: lhs, RHS: rhs})
			} else if rhs != nil {
				push(rhs)
			}
		case 21: // %
			rhs := pop()
			lhs := pop()
			if lhs != nil && rhs != nil {
				push(&image_formula_find.Modulus{LHS: lhs, RHS: rhs})
			} else if rhs != nil {
				push(rhs)
			}

		// Unary Ops (22-32)
		case 22: // Negate
			e := pop()
			if e != nil {
				push(&image_formula_find.Negate{Expr: e})
			}
		case 23: // Sin
			e := pop()
			if e != nil {
				push(&image_formula_find.SingleFunction{Name: "Sin", Expr: e})
			}
		case 24: // Cos
			e := pop()
			if e != nil {
				push(&image_formula_find.SingleFunction{Name: "Cos", Expr: e})
			}
		case 25: // Tan
			e := pop()
			if e != nil {
				push(&image_formula_find.SingleFunction{Name: "Tan", Expr: e})
			}
		case 26: // Asin
			e := pop()
			if e != nil {
				push(&image_formula_find.SingleFunction{Name: "Asin", Expr: e})
			}
		case 27: // Acos
			e := pop()
			if e != nil {
				push(&image_formula_find.SingleFunction{Name: "Acos", Expr: e})
			}
		case 28: // Atan
			e := pop()
			if e != nil {
				push(&image_formula_find.SingleFunction{Name: "Atan", Expr: e})
			}
		case 29: // Log
			e := pop()
			if e != nil {
				push(&image_formula_find.SingleFunction{Name: "Log", Expr: e})
			}
		case 30: // Exp
			e := pop()
			if e != nil {
				push(&image_formula_find.SingleFunction{Name: "Exp", Expr: e})
			}
		case 31: // Abs
			e := pop()
			if e != nil {
				push(&image_formula_find.SingleFunction{Name: "Abs", Expr: e})
			}
		case 32: // Sqrt
			e := pop()
			if e != nil {
				push(&image_formula_find.SingleFunction{Name: "Sqrt", Expr: e})
			}

		// Double Functions (33-34)
		case 33: // Min
			rhs := pop()
			lhs := pop()
			if lhs != nil && rhs != nil {
				push(&image_formula_find.DoubleFunction{Name: "Min", Expr1: lhs, Expr2: rhs})
			} else if rhs != nil {
				push(rhs)
			}
		case 34: // Max
			rhs := pop()
			lhs := pop()
			if lhs != nil && rhs != nil {
				push(&image_formula_find.DoubleFunction{Name: "Max", Expr1: lhs, Expr2: rhs})
			} else if rhs != nil {
				push(rhs)
			}
		case 35: // Sinh
			e := pop()
			if e != nil {
				push(&image_formula_find.SingleFunction{Name: "Sinh", Expr: e})
			}
		case 36: // Cosh
			e := pop()
			if e != nil {
				push(&image_formula_find.SingleFunction{Name: "Cosh", Expr: e})
			}
		case 37: // Tanh
			e := pop()
			if e != nil {
				push(&image_formula_find.SingleFunction{Name: "Tanh", Expr: e})
			}
		case 38: // Ceil
			e := pop()
			if e != nil {
				push(&image_formula_find.SingleFunction{Name: "Ceil", Expr: e})
			}
		case 39: // Floor
			e := pop()
			if e != nil {
				push(&image_formula_find.SingleFunction{Name: "Floor", Expr: e})
			}
		case 40: // Round
			e := pop()
			if e != nil {
				push(&image_formula_find.SingleFunction{Name: "Round", Expr: e})
			}
		case 41: // Atan2
			rhs := pop()
			lhs := pop()
			if lhs != nil && rhs != nil {
				push(&image_formula_find.DoubleFunction{Name: "Atan2", Expr1: lhs, Expr2: rhs})
			} else if rhs != nil {
				push(rhs)
			}
		case 42: // Hypot
			rhs := pop()
			lhs := pop()
			if lhs != nil && rhs != nil {
				push(&image_formula_find.DoubleFunction{Name: "Hypot", Expr1: lhs, Expr2: rhs})
			} else if rhs != nil {
				push(rhs)
			}
		case 43: // Dim
			rhs := pop()
			lhs := pop()
			if lhs != nil && rhs != nil {
				push(&image_formula_find.DoubleFunction{Name: "Dim", Expr1: lhs, Expr2: rhs})
			} else if rhs != nil {
				push(rhs)
			}

		default:
			// Map other chars to small random constants based on index
			// Indices not used: 11-15, 44-63.
			// Let's use them for more constants.
			v := float64(idx-44) / 5.0
			push(&image_formula_find.Const{Value: v})
		}
	}

	if len(stack) == 0 {
		return &image_formula_find.Const{Value: 0}
	}

	// Sum up the stack to use all generated parts
	var res image_formula_find.Expression = stack[0]
	for i := 1; i < len(stack); i++ {
		res = &image_formula_find.Plus{LHS: res, RHS: stack[i]}
	}
	return res
}

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
		e := (len(s) / p) * (i + 1)
		result += s[st:e]
	}
	return result
}

func Valid(dna string) bool {
	rf, bf, gf := ParseDNA(dna)
	return rf.HasVar("X") && rf.HasVar("Y") &&
		bf.HasVar("X") && bf.HasVar("Y") &&
		gf.HasVar("X") && gf.HasVar("Y")
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
		dna, ok := <-newDNA
		if !ok {
			break
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
