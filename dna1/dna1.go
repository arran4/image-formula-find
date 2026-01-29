package dna1

import (
	"github.com/agnivade/levenshtein"
	"image-formula-find"
	"math"
	"math/rand"
	"sort"
	"sync"
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

func ParseExpression(arg string) (string, image_formula_find.Expression) {
	if len(arg) == 0 {
		return "", &image_formula_find.Const{Value: 0}
	}
	c := arg[0]
	arg = arg[1:]
	i, ok := runeMapPos[rune(c)]
	if !ok {
		return MakeConst(arg, c)
	}

	switch i % 17 {
	case 1:
		lhs, rhs := Split2AndParse(arg)
		return "", image_formula_find.Plus{
			LHS: lhs,
			RHS: rhs,
		}
	case 2:
		lhs, rhs := Split2AndParse(arg)
		return "", image_formula_find.Subtract{
			LHS: lhs,
			RHS: rhs,
		}
	case 3:
		lhs, rhs := Split2AndParse(arg)
		return "", image_formula_find.Multiply{
			LHS: lhs,
			RHS: rhs,
		}
	case 4:
		lhs, rhs := Split2AndParse(arg)
		return "", image_formula_find.Divide{
			LHS: lhs,
			RHS: rhs,
		}
	case 5:
		lhs, rhs := Split2AndParse(arg)
		return "", image_formula_find.Modulus{
			LHS: lhs,
			RHS: rhs,
		}
	case 6:
		var farg rune = 'A'
		if len(arg) > 0 {
			farg = rune(arg[0])
			arg = arg[1:]
		}
		fi, _ := runeMapPos[farg]
		lhs, rhs := Split2AndParse(arg)
		return "", image_formula_find.DoubleFunction{
			Name:  image_formula_find.FunctionNames[fi%len(image_formula_find.FunctionNames)],
			Expr1: lhs,
			Expr2: rhs,
			Infix: false,
		}
	case 7:
		var farg rune = 'A'
		if len(arg) > 0 {
			farg = rune(arg[0])
			arg = arg[1:]
		}
		fi, _ := runeMapPos[farg]
		arg, expr := ParseExpression(arg)
		return arg, image_formula_find.SingleFunction{
			Name: image_formula_find.FunctionNames[int(fi)%len(image_formula_find.FunctionNames)],
			Expr: expr,
		}
	case 8:
		arg, expr := ParseExpression(arg)
		return arg, image_formula_find.Negate{
			Expr: expr,
		}
	case 9:
		lhs, rhs := Split2AndParse(arg)
		return "", image_formula_find.Power{
			LHS: lhs,
			RHS: rhs,
		}
	case 10:
		return arg, image_formula_find.Var{
			Var: "X",
		}
	case 11:
		return arg, image_formula_find.Var{
			Var: "Y",
		}
	case 0:
		return ParseConstWithExponent(arg, 0)
	case 12:
		return ParseConstWithExponent(arg, 1)
	case 13:
		return ParseConstWithExponent(arg, 2)
	case 14:
		return ParseConstWithExponent(arg, -1)
	case 15:
		return ParseConstWithExponent(arg, -2)
	case 16:
		return MakeConst(arg, c)
	default:
		return MakeConst(arg, c)
	}
}

func MakeConst(arg string, c uint8) (string, image_formula_find.Expression) {
	i := runeMapPos[rune(c)]
	return arg, &image_formula_find.Const{Value: float64(i)}
}

// Deprecated: Use ParseConstWithExponent
func ParseConst(arg string, c uint8) (string, image_formula_find.Expression) {
	var exponent int
	switch c {
	case chars[0]: // A
		exponent = 0
	case chars[16]: // Q
		exponent = 1
	case chars[32]: // g
		exponent = 2
	case chars[48]: // w
		exponent = -1
	case chars[63]: // /
		exponent = -2
	default:
		exponent = 0
	}
	return ParseConstWithExponent(arg, exponent)
}

func ParseConstWithExponent(arg string, exponent int) (string, image_formula_find.Expression) {
	m := math.Pow10(exponent)
	r := 0.0
	for i := 0; i < 2; i++ {
		if len(arg) == 0 {
			break
		}
		v, ok := runeMapPos[rune(arg[0])]
		if !ok {
			arg = arg[1:]
			continue
		}
		r = r*64 + float64(v)
		arg = arg[1:]
	}
	return arg, &image_formula_find.Const{Value: r * m}
}

func ParseExpressionAll(arg string) image_formula_find.Expression {
	if len(arg) == 0 {
		return &image_formula_find.Const{Value: 0}
	}
	arg, result := ParseExpression(arg)
	if len(arg) > 0 {
		result = image_formula_find.Plus{
			LHS: result,
			RHS: ParseExpressionAll(arg),
		}
	}
	return result
}

func ParseFunction(arg string) *image_formula_find.Function {
	lhs, rhs := Split2AndParse(arg)
	return &image_formula_find.Function{
		Equals: &image_formula_find.Equals{
			LHS: lhs,
			RHS: rhs,
		},
	}
}

func Split2AndParse(arg string) (image_formula_find.Expression, image_formula_find.Expression) {
	lhsStr, rhsStr := SplitString2(arg)
	lhs := ParseExpressionAll(lhsStr)
	rhs := ParseExpressionAll(rhsStr)
	return lhs, rhs
}

func SplitString2(arg string) (string, string) {
	lhsStr := ""
	rhsStr := ""
	for len(arg) > 0 {
		c := arg[0]
		arg = arg[1:]
		i, ok := runeMapPos[rune(c)]
		if !ok {
			continue
		}
		width := float64(len(arg) - 2)
		if width < 0 {
			break
		}
		incs := width / 64.0
		lhsStr = arg[0 : int(incs*float64(i))+1]
		rhsStr = arg[int(incs*float64(i))+1:]
		return lhsStr, rhsStr
	}
	return lhsStr, rhsStr
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
		e := (len(s)/p)*(i+1) - 1
		result += s[st:e]
	}
	return result
}

func ParseDNA(dna string) (*image_formula_find.Function, *image_formula_find.Function, *image_formula_find.Function) {
	rd, bd, gd := SplitString3(dna)
	rf := ParseFunction(rd)
	bf := ParseFunction(bd)
	gf := ParseFunction(gd)
	return rf, bf, gf
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

	//for i := 0; i < len(lastGeneration)*len(lastGeneration); i++ {
	//	p1 := lastGeneration[i%len(lastGeneration)]
	//	p2 := lastGeneration[i/len(lastGeneration)]
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
