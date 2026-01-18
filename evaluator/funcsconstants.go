package evaluator

import (
	"fmt"
	"math"
	"math/big"
)

var ErrInvalidArgCount = fmt.Errorf("invalid number of arguments")
var ErrIntegerExpected = fmt.Errorf("integer argument expected")
var ErrNaturalExpected = fmt.Errorf("natural number argument expected")

var Constants = map[string]Number{
	"pi": NewFloat(big.NewFloat(math.Pi).SetPrec(DefaultPrec)),
	"e":  NewFloat(big.NewFloat(math.E).SetPrec(DefaultPrec)),
}

type Function func(args []Number, degrees bool) (Number, error)

func fact(n Number) (Number, error) {
	if !n.IsInt() {
		return Zero(), ErrNaturalExpected
	}
	if n.Int().Sign() < 0 {
		return Zero(), ErrNaturalExpected
	}

	r := big.NewInt(1)
	i := big.NewInt(2)
	for i.Cmp(n.Int()) <= 0 {
		r.Mul(r, i)
		i.Add(i, big.NewInt(1))
	}
	return NewInt(r), nil
}

func unaryFloat(fn func(float64) float64) Function {
	return func(args []Number, degrees bool) (Number, error) {
		if len(args) != 1 {
			return Zero(), ErrInvalidArgCount
		}
		x, _ := args[0].Float().Float64()
		if degrees {
			x *= math.Pi / 180
		}
		return NewFloat(
			new(big.Float).SetPrec(DefaultPrec).SetFloat64(fn(x)),
		), nil
	}
}

func binaryFloat(fn func(float64, float64) float64) Function {
	return func(args []Number, degrees bool) (Number, error) {
		if len(args) != 2 {
			return Zero(), ErrInvalidArgCount
		}
		f641, _ := args[0].Float().Float64()
		f642, _ := args[1].Float().Float64()
		return NewFloat(
			new(big.Float).SetPrec(DefaultPrec).SetFloat64(
				fn(f641, f642),
			),
		), nil
	}
}

var Functions = map[string]Function{
	"sin":  unaryFloat(math.Sin),
	"cos":  unaryFloat(math.Cos),
	"tan":  unaryFloat(math.Tan),
	"cot":  unaryFloat(func(x float64) float64 { return 1 / math.Tan(x) }),

	"asin": unaryFloat(math.Asin),
	"acos": unaryFloat(math.Acos),
	"atan": unaryFloat(math.Atan),
	"atan2": binaryFloat(math.Atan2),

	"sinh":  unaryFloat(math.Sinh),
	"cosh":  unaryFloat(math.Cosh),
	"tanh":  unaryFloat(math.Tanh),
	"asinh": unaryFloat(math.Asinh),
	"acosh": unaryFloat(math.Acosh),
	"atanh": unaryFloat(math.Atanh),

	"ln":    unaryFloat(math.Log),
	"log":   binaryFloat(func(a, b float64) float64 { return math.Log(b) / math.Log(a) }),
	"log2":  unaryFloat(math.Log2),
	"log10": unaryFloat(math.Log10),
	"exp":   unaryFloat(math.Exp),

	"sqrt": unaryFloat(math.Sqrt),
	"cbrt": unaryFloat(math.Cbrt),
	"pow":  binaryFloat(math.Pow),
	"root": binaryFloat(func(a, b float64) float64 { return math.Pow(a, 1/b) }),
	"hypot": binaryFloat(math.Hypot),

	"ceil":  unaryFloat(math.Ceil),
	"floor": unaryFloat(math.Floor),
	"round": unaryFloat(math.Round),
	"trunc": unaryFloat(math.Trunc),

	"abs": unaryFloat(math.Abs),
	"min": binaryFloat(math.Min),
	"max": binaryFloat(math.Max),
	"mod": binaryFloat(math.Mod),

	"fact": func(args []Number, _ bool) (Number, error) {
		if len(args) != 1 {
			return Zero(), ErrInvalidArgCount
		}
		return fact(args[0])
	},
}
