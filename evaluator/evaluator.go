package evaluator

import (
	"fmt"
	"math"
	"math/big"

	"github.com/marzeq/mexe/parser"
)

func Evaluate(expr parser.Node, degreesMode bool) (Number, error) {
	switch n := expr.(type) {

	case *parser.ConstantNode:
		if v, err := ParseNumber(n.Name); err == nil {
			return v, nil
		}
		if c, ok := Constants[n.Name]; ok {
			return c, nil
		}
		return Zero(), fmt.Errorf("unknown constant: %s", n.Name)

	case *parser.UnaryOpNode:
		val, err := Evaluate(n.Child, degreesMode)
		if err != nil {
			return Zero(), err
		}
		switch n.OpType {
		case parser.UnaryOpTypeNegate:
			if val.IsInt() {
				return NewInt(new(big.Int).Neg(val.Int())), nil
			}
			return NewFloat(new(big.Float).Neg(val.Float())), nil

		case parser.UnaryOpTypeAbs:
			if val.IsInt() {
				return NewInt(new(big.Int).Abs(val.Int())), nil
			}
			return NewFloat(new(big.Float).Abs(val.Float())), nil

		case parser.UnaryOpTypeFact:
			return fact(val)

		default:
			return Zero(), fmt.Errorf("unknown unary operator")
		}

	case *parser.BinaryOpNode:
		lhs, err := Evaluate(n.Left, degreesMode)
		if err != nil {
			return Zero(), err
		}
		rhs, err := Evaluate(n.Right, degreesMode)
		if err != nil {
			return Zero(), err
		}

		switch n.OpType {
		case parser.BinaryOpTypeAdd:
			a, b := promote(lhs, rhs)
			if a.IsInt() {
				return NewInt(new(big.Int).Add(a.Int(), b.Int())), nil
			}
			return NewFloat(new(big.Float).Add(a.Float(), b.Float())), nil

		case parser.BinaryOpTypeSubtract:
			a, b := promote(lhs, rhs)
			if a.IsInt() {
				return NewInt(new(big.Int).Sub(a.Int(), b.Int())), nil
			}
			return NewFloat(new(big.Float).Sub(a.Float(), b.Float())), nil

		case parser.BinaryOpTypeMultiply:
			a, b := promote(lhs, rhs)
			if a.IsInt() {
				return NewInt(new(big.Int).Mul(a.Int(), b.Int())), nil
			}
			return NewFloat(new(big.Float).Mul(a.Float(), b.Float())), nil

		case parser.BinaryOpTypeDivide:
			if rhs.IsInt() && rhs.Int().Sign() == 0 {
				return Zero(), fmt.Errorf("division by zero")
			}
			return NewFloat(new(big.Float).Quo(lhs.Float(), rhs.Float())), nil

		case parser.BinaryOpTypeExponent:
			if lhs.IsInt() && rhs.IsInt() && rhs.Int().Sign() >= 0 {
				return NewInt(new(big.Int).Exp(lhs.Int(), rhs.Int(), nil)), nil
			}
			lhsf, _ := lhs.Float().Float64()
			rhsf, _ := rhs.Float().Float64()
			f := math.Pow(
				lhsf,
				rhsf,
			)
			return NewFloat(big.NewFloat(f).SetPrec(DefaultPrec)), nil

		case parser.BinaryOpTypeMod:
			if !lhs.IsInt() || !rhs.IsInt() {
				return Zero(), ErrIntegerExpected
			}
			return NewInt(new(big.Int).Mod(lhs.Int(), rhs.Int())), nil

		default:
			return Zero(), fmt.Errorf("unknown binary operator")
		}

	case *parser.FuncApplicationNode:
		fn, ok := Functions[n.FuncName]
		if !ok {
			return Zero(), fmt.Errorf("unknown function: %s", n.FuncName)
		}

		args := make([]Number, len(n.Args))
		for i, a := range n.Args {
			v, err := Evaluate(a, degreesMode)
			if err != nil {
				return Zero(), err
			}
			args[i] = v
		}

		return fn(args, degreesMode)

	default:
		return Zero(), fmt.Errorf("unknown node type")
	}
}
