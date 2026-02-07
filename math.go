package image_formula_find

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
)

var (
	SingleFunctions map[string]SingleFunctionDef
	DoubleFunctions map[string]DoubleFunctionDef
	FunctionNames   []string
)

type SingleFunctionDef func(float64) float64
type DoubleFunctionDef func(float64, float64) float64

type State struct {
	X, Y                            float64
	T                               int
	AccessedX, AccessedY, AccessedT bool
}

func (rs *State) CurX() float64 {
	rs.AccessedX = true
	return rs.X
}

func (rs *State) CurY() float64 {
	rs.AccessedY = true
	return rs.Y
}

func (rs *State) CurT() int {
	rs.AccessedT = true
	return rs.T
}

type Expression interface {
	Evaluate(state *State) float64
	String() string
	Depth() int
	Simplify() Expression
	HasVar(vs string) bool
}

type Function struct {
	Equals *Equals
}

var statePool = sync.Pool{
	New: func() interface{} {
		return &State{}
	},
}

func (v Function) Evaluate(X, Y float64, T int) (weight float64, TUsed bool, err error) {
	if v.Equals == nil {
		return 0, false, errors.New("no such formula")
	}

	state := statePool.Get().(*State)
	state.X = X
	state.Y = Y
	state.T = T
	state.AccessedX = false
	state.AccessedY = false
	state.AccessedT = false

	weight = v.Equals.Evaluate(state)
	TUsed = state.AccessedT
	statePool.Put(state)
	return
}

func (v Function) String() string {
	return v.Equals.String()
}

func (v Function) Simplify() *Function {
	e := v.Equals.Simplify().(*Equals)
	v.Equals = e
	return &v
}

func (v Function) HasVar(vs string) bool {
	return v.Equals.HasVar(vs)
}

type Equals struct {
	LHS Expression
	RHS Expression
}

func (v Equals) Evaluate(state *State) float64 {
	if v.LHS == nil {
		return v.RHS.Evaluate(state)
	}
	return v.RHS.Evaluate(state) - v.LHS.Evaluate(state)
}

func (v Equals) Depth() int {
	if v.LHS == nil {
		return v.RHS.Depth()
	}
	l, r := v.LHS.Depth(), v.RHS.Depth()
	if l > r {
		return l + 1
	}
	return r + 1
}

func (v Equals) String() string {
	if v.LHS == nil {
		return v.RHS.String()
	}
	return fmt.Sprintf("%s = %s", v.LHS.String(), v.RHS.String())
}

func (v Equals) Simplify() Expression {
	v.RHS = removeBrackets(v.RHS.Simplify())
	if v.LHS != nil {
		v.LHS = removeBrackets(v.LHS.Simplify())
	}
	return &v
}

func (v Equals) HasVar(vs string) bool {
	if v.LHS == nil {
		return v.RHS.HasVar(vs)
	}
	return v.RHS.HasVar(vs) || v.LHS.HasVar(vs)
}

func removeBrackets(e Expression) Expression {
	switch e := e.(type) {
	case *Brackets:
		return removeBrackets(e.Expr)
	}
	return e
}

type Var struct {
	Var string
}

func (v Var) HasVar(vs string) bool {
	return v.Var == vs
}

func (v Var) Evaluate(state *State) float64 {
	if len(v.Var) == 1 {
		if v.Var[0] == 'x' || v.Var[0] == 'X' {
			return float64(state.CurX())
		}
		if v.Var[0] == 'y' || v.Var[0] == 'Y' {
			return float64(state.CurY())
		}
		if v.Var[0] == 't' || v.Var[0] == 'T' {
			return float64(state.CurT())
		}
	}
	switch strings.ToUpper(v.Var) {
	case "X":
		return float64(state.CurX())
	case "Y":
		return float64(state.CurY())
	case "T":
		return float64(state.CurT())
	default:
		return 0
	}
}

func (v Var) Depth() int {
	return 1
}

func (v Var) String() string {
	return v.Var
}

func (v Var) Simplify() Expression {
	return &v
}

type Const struct {
	Value float64
}

func (c Const) HasVar(vs string) bool {
	return false
}

func (c Const) Evaluate(state *State) float64 {
	return c.Value
}

func (v Const) String() string {
	return fmt.Sprintf("%g", v.Value)
}

func (v Const) Simplify() Expression {
	return &v
}

func (v Const) Depth() int {
	return 1
}

type Plus struct {
	LHS Expression
	RHS Expression
}

func (v Plus) HasVar(vs string) bool {
	return v.RHS.HasVar(vs) || v.LHS.HasVar(vs)
}

func (v Plus) Evaluate(state *State) float64 {
	return v.RHS.Evaluate(state) + v.LHS.Evaluate(state)
}

func (v Plus) String() string {
	return fmt.Sprintf("%s + %s", v.LHS.String(), v.RHS.String())
}

func (v Plus) Simplify() Expression {
	v.RHS = v.RHS.Simplify()
	v.LHS = v.LHS.Simplify()
	return &v
}

func (v Plus) Depth() int {
	l, r := v.LHS.Depth(), v.RHS.Depth()
	if l > r {
		return l + 1
	}
	return r + 1
}

type Subtract struct {
	LHS Expression
	RHS Expression
}

func (v Subtract) HasVar(vs string) bool {
	return v.RHS.HasVar(vs) || v.LHS.HasVar(vs)
}

func (v Subtract) Evaluate(state *State) float64 {
	return v.RHS.Evaluate(state) - v.LHS.Evaluate(state)
}

func (v Subtract) String() string {
	return fmt.Sprintf("%s - %s", v.LHS.String(), v.RHS.String())
}

func (v Subtract) Simplify() Expression {
	v.RHS = v.RHS.Simplify()
	v.LHS = v.LHS.Simplify()
	return &v
}

func (v Subtract) Depth() int {
	l, r := v.LHS.Depth(), v.RHS.Depth()
	if l > r {
		return l + 1
	}
	return r + 1
}

type Multiply struct {
	LHS Expression
	RHS Expression
}

func (v Multiply) HasVar(vs string) bool {
	return v.RHS.HasVar(vs) || v.LHS.HasVar(vs)
}

func (v Multiply) Evaluate(state *State) float64 {
	return v.RHS.Evaluate(state) * v.LHS.Evaluate(state)
}

func (v Multiply) String() string {
	return fmt.Sprintf("%s * %s", v.LHS.String(), v.RHS.String())
}

func (v Multiply) Simplify() Expression {
	v.RHS = v.RHS.Simplify()
	v.LHS = v.LHS.Simplify()
	return &v
}

func (v Multiply) Depth() int {
	l, r := v.LHS.Depth(), v.RHS.Depth()
	if l > r {
		return l + 1
	}
	return r + 1
}

type Divide struct {
	LHS Expression
	RHS Expression
}

func (v Divide) HasVar(vs string) bool {
	return v.RHS.HasVar(vs) || v.LHS.HasVar(vs)
}

func (v Divide) Evaluate(state *State) float64 {
	return v.RHS.Evaluate(state) / v.LHS.Evaluate(state)
}

func (v Divide) String() string {
	return fmt.Sprintf("%s / %s", v.LHS.String(), v.RHS.String())
}

func (v Divide) Simplify() Expression {
	v.RHS = v.RHS.Simplify()
	v.LHS = v.LHS.Simplify()
	return &v
}

func (v Divide) Depth() int {
	l, r := v.LHS.Depth(), v.RHS.Depth()
	if l > r {
		return l + 1
	}
	return r + 1
}

type Power struct {
	LHS Expression
	RHS Expression
}

func (v Power) HasVar(vs string) bool {
	return v.RHS.HasVar(vs) || v.LHS.HasVar(vs)
}

func (v Power) Evaluate(state *State) float64 {
	return math.Pow(v.LHS.Evaluate(state), v.RHS.Evaluate(state))
}

func (v Power) String() string {
	return fmt.Sprintf("%s ^ %s", v.LHS.String(), v.RHS.String())
}

func (v Power) Simplify() Expression {
	v.RHS = v.RHS.Simplify()
	v.LHS = v.LHS.Simplify()
	return &v
}

func (v Power) Depth() int {
	l, r := v.LHS.Depth(), v.RHS.Depth()
	if l > r {
		return l + 1
	}
	return r + 1
}

type Modulus struct {
	LHS Expression
	RHS Expression
}

func (v Modulus) HasVar(vs string) bool {
	return v.RHS.HasVar(vs) || v.LHS.HasVar(vs)
}

func (v Modulus) Evaluate(state *State) float64 {
	return math.Mod(v.LHS.Evaluate(state), v.RHS.Evaluate(state))
}

func (v Modulus) String() string {
	return fmt.Sprintf("%s %% %s", v.LHS.String(), v.RHS.String())
}

func (v Modulus) Simplify() Expression {
	v.RHS = v.RHS.Simplify()
	v.LHS = v.LHS.Simplify()
	return &v
}

func (v Modulus) Depth() int {
	l, r := v.LHS.Depth(), v.RHS.Depth()
	if l > r {
		return l + 1
	}
	return r + 1
}

type Negate struct {
	Expr Expression
}

func (v Negate) HasVar(vs string) bool {
	return v.Expr.HasVar(vs)
}

func (v Negate) Evaluate(state *State) float64 {
	return -v.Expr.Evaluate(state)
}

func (v Negate) String() string {
	return fmt.Sprintf("-%s", v.Expr.String())
}

func (v Negate) Simplify() Expression {
	switch child := v.Expr.(type) {
	case *Negate:
		return child.Expr.Simplify()
	case *Brackets:
		switch childChild := child.Expr.(type) {
		case *Negate:
			return childChild.Expr.Simplify()
		case *Const, *Var:
			v.Expr = childChild
		}
	}
	v.Expr = v.Expr.Simplify()
	return &v
}

func (v Negate) Depth() int {
	return v.Expr.Depth() + 1
}

type Brackets struct {
	Expr Expression
}

func (v Brackets) HasVar(vs string) bool {
	return v.Expr.HasVar(vs)
}

func (v Brackets) Evaluate(state *State) float64 {
	return v.Expr.Evaluate(state)
}

func (v Brackets) String() string {
	return fmt.Sprintf("(%s)", v.Expr.String())
}

func (v Brackets) Simplify() Expression {
	switch next := v.Expr.(type) {
	case *Brackets:
		return next.Expr.Simplify()
	case *Const:
		return next.Simplify()
	case *Var:
		return next.Simplify()
	}
	v.Expr = v.Expr.Simplify()
	return &v
}

func (v Brackets) Depth() int {
	return v.Expr.Depth() + 1
}

type SingleFunction struct {
	Name string
	Expr Expression
}

func (v SingleFunction) HasVar(vs string) bool {
	return v.Expr.HasVar(vs)
}

func (v SingleFunction) Evaluate(state *State) float64 {
	var r = v.Expr.Evaluate(state)
	if f, ok := SingleFunctions[strings.ToUpper(v.Name)]; ok {
		r = f(r)
	}
	return r
}

func (v SingleFunction) String() string {
	return fmt.Sprintf("%s(%s)", v.Name, v.Expr.String())
}

func (v SingleFunction) Simplify() Expression {
	v.Expr = v.Expr.Simplify()
	return &v
}

func (v SingleFunction) Depth() int {
	return v.Expr.Depth() + 1
}

type DoubleFunction struct {
	Name  string
	Expr1 Expression
	Expr2 Expression
	Infix bool
}

func (v DoubleFunction) HasVar(vs string) bool {
	return v.Expr1.HasVar(vs) || v.Expr2.HasVar(vs)
}

func (v DoubleFunction) Evaluate(state *State) float64 {
	var r1 = v.Expr1.Evaluate(state)
	var r2 = v.Expr2.Evaluate(state)
	if f, ok := DoubleFunctions[strings.ToUpper(v.Name)]; ok {
		r1 = f(r1, r2)
	}
	return r1
}

func (v DoubleFunction) String() string {
	if v.Infix {
		return fmt.Sprintf("%s %s %s", v.Expr1.String(), v.Name, v.Expr2.String())
	} else {
		return fmt.Sprintf("%s(%s, %s)", v.Name, v.Expr1.String(), v.Expr2.String())
	}
}

func (v DoubleFunction) Simplify() Expression {
	v.Expr1 = v.Expr1.Simplify()
	v.Expr2 = v.Expr2.Simplify()
	return &v
}

func (v DoubleFunction) Depth() int {
	l, r := v.Expr1.Depth(), v.Expr2.Depth()
	if l > r {
		return l + 1
	}
	return r + 1
}

func init() {
	SingleFunctions = map[string]SingleFunctionDef{}
	DoubleFunctions = map[string]DoubleFunctionDef{}
	FunctionNames = []string{}
	for name, f := range map[string]interface{}{
		"Abs":             math.Abs,
		"Acos":            math.Acos,
		"Acosh":           math.Acosh,
		"Asin":            math.Asin,
		"Asinh":           math.Asinh,
		"Atan":            math.Atan,
		"Atan2":           math.Atan2,
		"Atanh":           math.Atanh,
		"Cbrt":            math.Cbrt,
		"Ceil":            math.Ceil,
		"Copysign":        math.Copysign,
		"Cos":             math.Cos,
		"Cosh":            math.Cosh,
		"Dim":             math.Dim,
		"Erf":             math.Erf,
		"Erfc":            math.Erfc,
		"Erfcinv":         math.Erfcinv,
		"Erfinv":          math.Erfinv,
		"Exp":             math.Exp,
		"Exp2":            math.Exp2,
		"Expm1":           math.Expm1,
		"Float32bits":     math.Float32bits,
		"Float32frombits": math.Float32frombits,
		"Float64bits":     math.Float64bits,
		"Float64frombits": math.Float64frombits,
		"Floor":           math.Floor,
		"Frexp":           math.Frexp,
		"Gamma":           math.Gamma,
		"Hypot":           math.Hypot,
		"Ilogb":           math.Ilogb,
		"Inf":             math.Inf,
		"IsInf":           math.IsInf,
		"IsNaN":           math.IsNaN,
		"J0":              math.J0,
		"J1":              math.J1,
		"Jn":              math.Jn,
		"Ldexp":           math.Ldexp,
		"Lgamma":          math.Lgamma,
		"Log":             math.Log,
		"Log10":           math.Log10,
		"Log1p":           math.Log1p,
		"Log2":            math.Log2,
		"Logb":            math.Logb,
		"Max":             math.Max,
		"Min":             math.Min,
		"Mod":             math.Mod,
		"Modf":            math.Modf,
		"NaN":             math.NaN,
		"Nextafter":       math.Nextafter,
		"Nextafter32":     math.Nextafter32,
		"Pow":             math.Pow,
		"Pow10":           math.Pow10,
		"Remainder":       math.Remainder,
		"Round":           math.Round,
		"RoundToEven":     math.RoundToEven,
		"Signbit":         math.Signbit,
		"Sin":             math.Sin,
		"Sincos":          math.Sincos,
		"Sinh":            math.Sinh,
		"Sqrt":            math.Sqrt,
		"Tan":             math.Tan,
		"Tanh":            math.Tanh,
		"Trunc":           math.Trunc,
		"Y0":              math.Y0,
		"Y1":              math.Y1,
		"Yn":              math.Yn,
	} {
		switch f := f.(type) {
		case func(float64) float64:
			SingleFunctions[strings.ToUpper(name)] = f
		case func(float64, float64) float64:
			DoubleFunctions[strings.ToUpper(name)] = f
		case func(int, float64) float64:
			DoubleFunctions[strings.ToUpper(name)] = func(f1 float64, f2 float64) float64 {
				return f(int(f1), f2)
			}
		case func(float64, int) float64:
			DoubleFunctions[strings.ToUpper(name)] = func(f1 float64, f2 float64) float64 {
				return f(f1, int(f2))
			}
		case func(int) float64:
			SingleFunctions[strings.ToUpper(name)] = func(f1 float64) float64 {
				return f(int(f1))
			}
		case func(float64) int:
			SingleFunctions[strings.ToUpper(name)] = func(f1 float64) float64 {
				return float64(f(f1))
			}
		default:
			continue
		}
		FunctionNames = append(FunctionNames, name)
	}
}

func ParseFunction(arg string) *Function {
	yyResult = nil
	if r := yyParse(NewCalcLexer(arg)); r != 0 {
		log.Panic("Invalid formula: ", arg, " Left with ", yyResult)
	}
	return yyResult
}
