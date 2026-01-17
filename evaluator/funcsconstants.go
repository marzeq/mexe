package evaluator

import (
	"fmt"
	"math"
)

type Constant float64

type AngleMode int

const (
	AngleNone AngleMode = iota
	AngleIn
	AngleOut
)

type Function struct {
	Fn    func(args ...float64) (float64, error)
	Angle AngleMode
}

var ErrInvalidArgCount = fmt.Errorf("invalid number of arguments")
var ErrIntegerExpected = fmt.Errorf("integer argument expected")
var ErrNaturalExpected = fmt.Errorf("natural number argument expected")

var Constants = map[string]Constant{
	"pi": Constant(math.Pi),
	"e":  Constant(math.E),
}

func unary(fn func(float64) float64) func(args ...float64) (float64, error) {
	return func(args ...float64) (float64, error) {
		if len(args) != 1 {
			return 0, ErrInvalidArgCount
		}
		return fn(args[0]), nil
	}
}

func binary(fn func(float64, float64) float64) func(args ...float64) (float64, error) {
	return func(args ...float64) (float64, error) {
		if len(args) != 2 {
			return 0, ErrInvalidArgCount
		}
		return fn(args[0], args[1]), nil
	}
}

func fact(n float64) (float64, error) {
	if n < 0 || n != math.Trunc(n) {
		return 0, ErrNaturalExpected
	}
	r := 1.0
	for i := 2.0; i <= n; i++ {
		r *= i
	}
	return r, nil
}

var Functions = map[string]Function{
	"sin": {Fn: unary(math.Sin), Angle: AngleIn},
	"cos": {Fn: unary(math.Cos), Angle: AngleIn},
	"tan": {Fn: unary(math.Tan), Angle: AngleIn},
	"cot": {Fn: func(args ...float64) (float64, error) {
		if len(args) != 1 {
			return 0, ErrInvalidArgCount
		}
		return 1 / math.Tan(args[0]), nil
	}, Angle: AngleIn},

	"asin":    {Fn: unary(math.Asin), Angle: AngleOut},
	"arcsin":  {Fn: unary(math.Asin), Angle: AngleOut},
	"acos":    {Fn: unary(math.Acos), Angle: AngleOut},
	"arccos":  {Fn: unary(math.Acos), Angle: AngleOut},
	"atan":    {Fn: unary(math.Atan), Angle: AngleOut},
	"arctan":  {Fn: unary(math.Atan), Angle: AngleOut},
	"atan2":   {Fn: binary(math.Atan2), Angle: AngleOut},
	"arctan2": {Fn: binary(math.Atan2), Angle: AngleOut},

	"sinh":    {Fn: unary(math.Sinh)},
	"cosh":    {Fn: unary(math.Cosh)},
	"tanh":    {Fn: unary(math.Tanh)},
	"asinh":   {Fn: unary(math.Asinh)},
	"arcsinh": {Fn: unary(math.Asinh)},
	"acosh":   {Fn: unary(math.Acosh)},
	"arccosh": {Fn: unary(math.Acosh)},
	"atanh":   {Fn: unary(math.Atanh)},
	"arctanh": {Fn: unary(math.Atanh)},

	"ln": {Fn: unary(math.Log)},
	"log": {Fn: func(args ...float64) (float64, error) {
		if len(args) != 2 {
			return 0, ErrInvalidArgCount
		}
		return math.Log(args[1]) / math.Log(args[0]), nil
	}},
	"log2":  {Fn: unary(math.Log2)},
	"log10": {Fn: unary(math.Log10)},
	"exp":   {Fn: unary(math.Exp)},

	"sqrt": {Fn: unary(math.Sqrt)},
	"cbrt": {Fn: unary(math.Cbrt)},
	"pow":  {Fn: binary(math.Pow)},
	"root": {Fn: func(args ...float64) (float64, error) {
		if len(args) != 2 {
			return 0, ErrInvalidArgCount
		}
		return math.Pow(args[0], 1/args[1]), nil
	}},
	"hypot": {Fn: binary(math.Hypot)},

	"ceil":  {Fn: unary(math.Ceil)},
	"floor": {Fn: unary(math.Floor)},
	"round": {Fn: unary(math.Round)},
	"trunc": {Fn: unary(math.Trunc)},
	"fract": {Fn: func(args ...float64) (float64, error) {
		if len(args) != 1 {
			return 0, ErrInvalidArgCount
		}
		return args[0] - math.Trunc(args[0]), nil
	}},

	"abs": {Fn: unary(math.Abs)},
	"sign": {Fn: func(args ...float64) (float64, error) {
		if len(args) != 1 {
			return 0, ErrInvalidArgCount
		}
		switch {
		case args[0] > 0:
			return 1, nil
		case args[0] < 0:
			return -1, nil
		default:
			return 0, nil
		}
	}},
	"min": {Fn: func(args ...float64) (float64, error) {
		if len(args) < 1 {
			return 0, ErrInvalidArgCount
		}
		m := args[0]
		for _, v := range args[1:] {
			if v < m {
				m = v
			}
		}
		return m, nil
	}},
	"max": {Fn: func(args ...float64) (float64, error) {
		if len(args) < 1 {
			return 0, ErrInvalidArgCount
		}
		m := args[0]
		for _, v := range args[1:] {
			if v > m {
				m = v
			}
		}
		return m, nil
	}},
	"clamp": {Fn: func(args ...float64) (float64, error) {
		if len(args) != 3 {
			return 0, ErrInvalidArgCount
		}
		v, lo, hi := args[0], args[1], args[2]
		if v < lo {
			return lo, nil
		}
		if v > hi {
			return hi, nil
		}
		return v, nil
	}},
	"mod": {Fn: binary(math.Mod)},
	"fact": {Fn: func(args ...float64) (float64, error) {
		if len(args) != 1 {
			return 0, ErrInvalidArgCount
		}
		return fact(args[0])
	}},
}
