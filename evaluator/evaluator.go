package evaluator

import (
	"fmt"
	"math"

	"github.com/marzeq/mexe/parser"
)

func Evaluate(expr parser.Node, degreesMode bool) (float64, error) {
	switch n := expr.(type) {

	case *parser.ConstantNode:
		if v, err := parseNumber(n.Name); err == nil {
			return v, nil
		}

		if c, ok := Constants[n.Name]; ok {
			return float64(c), nil
		}

		return 0, fmt.Errorf("unknown constant: %s", n.Name)

	case *parser.UnaryOpNode:
		val, err := Evaluate(n.Child, degreesMode)
		if err != nil {
			return 0, err
		}

		switch n.OpType {
		case parser.UnaryOpTypeNegate:
			return -val, nil
		case parser.UnaryOpTypeAbs:
			return math.Abs(val), nil
		case parser.UnaryOpTypeFact:
			return fact(val)
		default:
			return 0, fmt.Errorf("unknown unary operator")
		}

	case *parser.BinaryOpNode:
		lhs, err := Evaluate(n.Left, degreesMode)
		if err != nil {
			return 0, err
		}
		rhs, err := Evaluate(n.Right, degreesMode)
		if err != nil {
			return 0, err
		}

		switch n.OpType {
		case parser.BinaryOpTypeAdd:
			return lhs + rhs, nil
		case parser.BinaryOpTypeSubtract:
			return lhs - rhs, nil
		case parser.BinaryOpTypeMultiply:
			return lhs * rhs, nil
		case parser.BinaryOpTypeDivide:
			if rhs == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return lhs / rhs, nil
		case parser.BinaryOpTypeExponent:
			return math.Pow(lhs, rhs), nil
		case parser.BinaryOpTypeMod:
			return math.Mod(lhs, rhs), nil
		default:
			return 0, fmt.Errorf("unknown binary operator")
		}

	case *parser.FuncApplicationNode:
		fn, ok := Functions[n.FuncName]
		if !ok {
			return 0, fmt.Errorf("unknown function: %s", n.FuncName)
		}

		args := make([]float64, len(n.Args))
		for i, arg := range n.Args {
			v, err := Evaluate(arg, degreesMode)
			if err != nil {
				return 0, err
			}
			args[i] = v
		}

		if degreesMode && fn.Angle == AngleIn {
			for i := range args {
				args[i] = args[i] * math.Pi / 180
			}
		}

		result, err := fn.Fn(args...)
		if err != nil {
			return 0, err
		}

		if degreesMode && fn.Angle == AngleOut {
			result = result * 180 / math.Pi
		}

		return result, nil

	default:
		return 0, fmt.Errorf("unknown node type")
	}
}

func parseNumber(s string) (float64, error) {
	var v float64
	_, err := fmt.Sscan(s, &v)
	return v, err
}
