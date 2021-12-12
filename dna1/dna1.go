package dna1

import (
	"image-formula-find"
	"math"
	"math/rand"
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
	switch c {
	case chars[1]:
		lhs, rhs := Split2AndParse(arg)
		return "", image_formula_find.Plus{
			LHS: lhs,
			RHS: rhs,
		}
	case chars[2]:
		lhs, rhs := Split2AndParse(arg)
		return "", image_formula_find.Subtract{
			LHS: lhs,
			RHS: rhs,
		}
	case chars[3]:
		lhs, rhs := Split2AndParse(arg)
		return "", image_formula_find.Multiply{
			LHS: lhs,
			RHS: rhs,
		}
	case chars[4]:
		lhs, rhs := Split2AndParse(arg)
		return "", image_formula_find.Divide{
			LHS: lhs,
			RHS: rhs,
		}
	case chars[5]:
		lhs, rhs := Split2AndParse(arg)
		return "", image_formula_find.Modulus{
			LHS: lhs,
			RHS: rhs,
		}
	case chars[6]:
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
	case chars[7]:
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
	case chars[8]:
		arg, expr := ParseExpression(arg)
		return arg, image_formula_find.Negate{
			Expr: expr,
		}
	case chars[9]:
		lhs, rhs := Split2AndParse(arg)
		return "", image_formula_find.Power{
			LHS: lhs,
			RHS: rhs,
		}
	case chars[10]:
		return arg, image_formula_find.Var{
			Var: "X",
		}
	case chars[11]:
		return arg, image_formula_find.Var{
			Var: "Y",
		}

	case chars[0], chars[len(chars)/2], chars[len(chars)-1]:
		return ParseConst(arg, c)
	default:
		return MakeConst(arg, c)
	}
}

func MakeConst(arg string, c uint8) (string, image_formula_find.Expression) {
	i := runeMapPos[rune(c)]
	return arg, &image_formula_find.Const{Value: float64(i)}
}

func ParseConst(arg string, c uint8) (string, image_formula_find.Expression) {
	m := math.Pow10(int(c >> 6))
	r := 0.0
	for i := 0; i < 2; i++ {
		if len(arg) == 0 {
			break
		}
		v, ok := runeMapPos[rune(arg[0])]
		if !ok {
			continue
		}
		r = r*64 + float64(v)
	}
	return arg, &image_formula_find.Const{Value: r / m}
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
		for i2-i1 < 2 {
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
